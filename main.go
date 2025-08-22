// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics
var (
	queueLength = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "email_queue_length",
		Help: "Current number of jobs in the queue",
	})
	jobsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "email_jobs_processed_total",
		Help: "Total number of successfully processed jobs",
	})
	jobsFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "email_jobs_failed_total",
		Help: "Total number of failed jobs sent to dead letter queue",
	})
)

func init() {
	prometheus.MustRegister(queueLength, jobsProcessed, jobsFailed)
}

func main() {
	// Configurable via env vars
	queueSize := 10
	numWorkers := 3
	if qs := os.Getenv("QUEUE_SIZE"); qs != "" {
		if val, err := strconv.Atoi(qs); err == nil {
			queueSize = val
		}
	}
	if nw := os.Getenv("NUM_WORKERS"); nw != "" {
		if val, err := strconv.Atoi(nw); err == nil {
			numWorkers = val
		}
	}

	queue := NewQueue(queueSize)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		go WorkerWithMetrics(ctx, i, queue)
	}

	// HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/send-email", SendEmailHandler(queue))

	// Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Update queue length in background
	go func() {
		for {
			queueLength.Set(float64(len(queue.Queue)))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	go func() {
		log.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	// Stop accepting new requests
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()
	server.Shutdown(ctxTimeout)

	// Stop workers
	cancel()

	// Close queue channel to signal no more jobs
	close(queue.Queue)

	log.Println("Server gracefully stopped")
	log.Printf("Dead Letter Queue size: %d\n", len(queue.DeadLetter))
}

// worker.go
package main

import (
	"context"
	"log"
	"math/rand"
)

const MaxRetries = 3

// WorkerWithMetrics processes jobs and updates Prometheus metrics
func WorkerWithMetrics(ctx context.Context, id int, queue *EmailQueue) {
	rand.Seed(int64(id)) // seed for randomness per worker

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d shutting down\n", id)
			return
		case job, ok := <-queue.Queue:
			if !ok {
				log.Printf("Worker %d: queue closed, exiting\n", id)
				return
			}

			attempts := 0
			for {
				attempts++
				err := job.SimulateSend()
				if err == nil {
					log.Printf("Worker %d successfully sent email to %s\n", id, job.To)
					jobsProcessed.Inc() // increment success metric
					break
				}

				log.Printf("Worker %d failed to send email to %s (attempt %d)\n", id, job.To, attempts)

				if attempts >= MaxRetries {
					log.Printf("Worker %d: max retries reached for %s, sending to Dead Letter Queue\n", id, job.To)
					queue.AddToDeadLetter(job)
					jobsFailed.Inc() // increment failure metric
					break
				}
			}
		}
	}
}

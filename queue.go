// queue.go
package main

import "fmt"

// EmailQueue holds the job channel
type EmailQueue struct {
	Queue      chan EmailJob
	Size       int
	DeadLetter []EmailJob
}

// NewQueue initializes the queue with a given size
func NewQueue(size int) *EmailQueue {
	return &EmailQueue{
		Queue:      make(chan EmailJob, size),
		Size:       size,
		DeadLetter: []EmailJob{},
	}
}

// Enqueue adds a job to the queue, returns error if full
func (q *EmailQueue) Enqueue(job EmailJob) error {
	select {
	case q.Queue <- job:
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

// AddToDeadLetter stores a permanently failed job
func (q *EmailQueue) AddToDeadLetter(job EmailJob) {
	q.DeadLetter = append(q.DeadLetter, job)
	fmt.Printf("Job added to Dead Letter Queue: %v\n", job)
}

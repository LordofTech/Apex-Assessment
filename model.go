// model.go
package main

import (
	"errors"
	"math/rand"
	"regexp"
	"time"
)

// EmailJob represents an email job
type EmailJob struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Retries int    `json:"-"`
}

// Validate checks the email job for required fields and a simple email format
func (e *EmailJob) Validate() error {
	if e.To == "" || e.Subject == "" || e.Body == "" {
		return errors.New("all fields are required")
	}

	// simple email regex
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(e.To) {
		return errors.New("invalid email address")
	}

	return nil
}

// SimulateSend randomly fails to simulate email sending errors
func (e *EmailJob) SimulateSend() error {
	time.Sleep(1 * time.Second) // simulate sending delay

	// 20% chance to fail
	if rand.Intn(100) < 20 {
		return errors.New("failed to send email")
	}

	return nil
}

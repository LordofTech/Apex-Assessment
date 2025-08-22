// handler.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// SendEmailHandler handles POST /send-email
func SendEmailHandler(queue *EmailQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var job EmailJob

		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if err := job.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if err := queue.Enqueue(job); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		log.Printf("Enqueued job to %s\n", job.To)
	}
}

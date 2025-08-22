# **README.md**

````markdown
# Go Email Microservice

A simple Go microservice to accept email jobs over HTTP, queue them, and process asynchronously using workers. Includes retry logic, dead-letter queue, configurable workers, queue size, and Prometheus metrics.

---

## **Features**

- `POST /send-email` endpoint to enqueue emails
- Input validation (all fields required, simple email check)
- In-memory job queue
- Multiple concurrent workers
- Retry failed jobs up to 3 times
- Dead Letter Queue for permanently failed jobs
- Graceful shutdown (`SIGINT` / `SIGTERM`)
- Configurable workers & queue size via environment variables
- Prometheus metrics (`/metrics` endpoint)
- Dockerized for easy deployment

---

## **API Usage**

### **POST /send-email**

**Request:**

```bash
curl -X POST http://localhost:8080/send-email \
-H "Content-Type: application/json" \
-d '{"to":"user@example.com","subject":"Hello","body":"Welcome!"}'
````

**Response Codes:**

* `202 Accepted` – Job enqueued
* `422 Bad Request` – Invalid input
* `503 Service Unavailable` – Queue is full

---

## **Prometheus Metrics**

Available at:

```
http://localhost:8080/metrics
```

**Metrics exposed:**

* `email_queue_length` – Current number of jobs in the queue
* `email_jobs_processed_total` – Total successfully processed jobs
* `email_jobs_failed_total` – Total jobs sent to Dead Letter Queue

---

## **Configuration**

Set via **environment variables**:

| Variable      | Default | Description                  |
| ------------- | ------- | ---------------------------- |
| `QUEUE_SIZE`  | 10      | Maximum jobs in the queue    |
| `NUM_WORKERS` | 3       | Number of concurrent workers |

---

## **Run Locally**

Make sure you have Go installed (>= 1.25).

1. Clone the repo:

```bash
git clone <your-repo-url>
cd go-email-microservice
```

2. Run the service:

```bash
go run main.go handler.go worker.go queue.go model.go
```

3. Test API with curl:

```bash
curl -X POST http://localhost:8080/send-email \
-H "Content-Type: application/json" \
-d '{"to":"user@example.com","subject":"Hello","body":"Welcome!"}'
```

---

## **Run with Docker**

1. Build the Docker image:

```bash
docker build -t go-email-service .
```

2. Run the container:

```bash
docker run -p 8080:8080 -e QUEUE_SIZE=20 -e NUM_WORKERS=5 go-email-service
```

3. Test the API (same curl command as above).

4. Access metrics:

```
http://localhost:8080/metrics
```

5. Stop the container:

```bash
docker ps
docker stop <container_id>
```

---

## **Graceful Shutdown**

* Press `Ctrl+C` in terminal or send `SIGINT`/`SIGTERM`
* Stops accepting new requests
* Waits for all active workers to finish
* Logs the final size of the Dead Letter Queue

---

## **Project Structure**

```
go-email-microservice/
│
├── main.go           # Entry point, server setup, graceful shutdown
├── handler.go        # HTTP handler for /send-email
├── worker.go         # Worker logic, retries, metrics
├── queue.go          # In-memory queue and Dead Letter Queue
├── model.go          # Job model and validation
├── go.mod
├── go.sum
└── Dockerfile
```

---

## **Notes**

* Default queue size and workers can be overridden using environment variables.
* Prometheus metrics are updated in real-time.
* The service is ready for submission or production deployment.

```

---

This README includes:  

- **Local and Docker run instructions**  
- **API usage and curl examples**  
- **Metrics and config explanation**  
- **Project structure**  




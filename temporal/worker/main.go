package main

import (
	"log"
	"temporal/activities"
	"temporal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

/*
WORKER ROLE:

The worker is responsible for:
1. Polling the Temporal server for tasks
2. Executing workflows and activities
3. Reporting results back to the server

Key concepts:
- Workers are stateless and can be scaled horizontally
- Multiple workers can share the same task queue
- Workers should be kept running (long-lived processes)
- Each worker can handle multiple concurrent executions
*/

func main() {
	// Create Temporal client
	// This connects to the Temporal server
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233", // Temporal server address

		// Namespace: logical isolation of workflows
		// "default" is the default namespace
		Namespace: "default",

		// Logger: configure logging for debugging
		// Logger: log.NewStructuredLogger(slog.New(...)),
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create worker
	// Workers poll a specific task queue for work
	w := worker.New(c, "order-processing-queue", worker.Options{
		// MaxConcurrentActivityExecutionSize: Max parallel activities
		MaxConcurrentActivityExecutionSize: 10,

		// MaxConcurrentWorkflowTaskExecutionSize: Max parallel workflow tasks
		MaxConcurrentWorkflowTaskExecutionSize: 10,

		// Enable session workers for file activities, etc.
		// EnableSessionWorker: true,
	})

	// Register workflows
	// The worker needs to know about all workflows it can execute
	w.RegisterWorkflow(workflows.OrderWorkflow)
	w.RegisterWorkflow(workflows.PaymentWorkflow) // Child workflow must also be registered

	// Register activities
	// Activities are registered as methods on a struct
	// This allows for dependency injection and better testability
	orderActivities := &activities.OrderActivities{}
	w.RegisterActivity(orderActivities.ValidateOrder)
	w.RegisterActivity(orderActivities.ReserveInventory)
	w.RegisterActivity(orderActivities.SendNotification)
	w.RegisterActivity(orderActivities.CompensateInventory)

	paymentActivities := &activities.PaymentActivities{}
	w.RegisterActivity(paymentActivities.AuthorizePayment)
	w.RegisterActivity(paymentActivities.CapturePayment)
	w.RegisterActivity(paymentActivities.RefundPayment)
	w.RegisterActivity(paymentActivities.VoidAuthorization)

	// Alternative: Register all methods on a struct
	// w.RegisterActivity(orderActivities)
	// w.RegisterActivity(paymentActivities)

	log.Println("Worker starting...")
	log.Println("Listening on task queue: order-processing-queue")
	log.Println("Press Ctrl+C to stop the worker")

	// Start worker
	// This is a blocking call - the worker will run until stopped
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

	log.Println("Worker stopped")
}

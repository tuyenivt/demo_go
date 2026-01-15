package main

import (
	"context"
	"fmt"
	"log"
	"temporal/models"
	"temporal/workflows"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

/*
CLIENT ROLE:

The client is responsible for:
1. Starting new workflow executions
2. Querying workflow state
3. Signaling running workflows
4. Canceling workflows

This is typically where your application code interacts with Temporal.
*/

func main() {
	// Create Temporal client
	c, err := client.Dial(client.Options{
		HostPort:  "localhost:7233",
		Namespace: "default",
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create a sample order
	order := models.Order{
		OrderID:    fmt.Sprintf("order-%s", uuid.New().String()[:8]),
		CustomerID: "customer-123",
		Items: []models.OrderItem{
			{
				ProductID: "prod-001",
				Quantity:  2,
				Price:     29.99,
			},
			{
				ProductID: "prod-002",
				Quantity:  1,
				Price:     49.99,
			},
		},
		TotalAmount: 109.97,
		Status:      "pending",
	}

	fmt.Println("=== Starting Order Processing Workflow ===")
	fmt.Printf("Order ID: %s\n", order.OrderID)
	fmt.Printf("Customer ID: %s\n", order.CustomerID)
	fmt.Printf("Total Amount: $%.2f\n", order.TotalAmount)
	fmt.Println()

	// Configure workflow options
	workflowOptions := client.StartWorkflowOptions{
		// Unique ID for this workflow execution
		// If a workflow with this ID is already running, it will fail
		ID: fmt.Sprintf("order-workflow-%s", order.OrderID),

		// Task queue where the workflow will be executed
		// Must match the task queue that workers are polling
		TaskQueue: "order-processing-queue",

		// WORKFLOW-LEVEL RETRY POLICY
		// This is different from activity retry!
		// Applied when the ENTIRE workflow fails and needs to restart from beginning
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    2 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    1 * time.Minute,
			MaximumAttempts:    3, // Retry the entire workflow up to 3 times
		},

		// Workflow execution timeout: Maximum time for workflow to complete
		WorkflowExecutionTimeout: 10 * time.Minute,

		// Workflow run timeout: Maximum time for a single workflow run
		// (useful if workflow continues-as-new)
		WorkflowRunTimeout: 10 * time.Minute,

		// Workflow task timeout: Maximum time to process a single workflow task
		WorkflowTaskTimeout: 10 * time.Second,

		// Memo: Static metadata about the workflow (searchable)
		// Memo: map[string]interface{}{
		//     "customerID": order.CustomerID,
		//     "orderTotal": order.TotalAmount,
		// },

		// SearchAttributes: Custom indexed fields for querying
		// SearchAttributes: map[string]interface{}{
		//     "CustomKeywordField": "order",
		// },
	}

	// Start the workflow
	// This is asynchronous - it returns immediately
	workflowRun, err := c.ExecuteWorkflow(
		context.Background(),
		workflowOptions,
		workflows.OrderWorkflow, // Workflow function
		order,                   // Workflow arguments
	)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	fmt.Printf("Started workflow\n")
	fmt.Printf("  Workflow ID: %s\n", workflowRun.GetID())
	fmt.Printf("  Run ID: %s\n", workflowRun.GetRunID())
	fmt.Println()
	fmt.Println("View workflow execution in Temporal UI:")
	fmt.Printf("  http://localhost:8233/namespaces/default/workflows/%s/%s\n",
		workflowRun.GetID(), workflowRun.GetRunID())
	fmt.Println()

	// Option 1: Wait for workflow to complete (blocking)
	fmt.Println("Waiting for workflow to complete...")
	var result string
	err = workflowRun.Get(context.Background(), &result)
	if err != nil {
		log.Printf("Workflow execution failed: %v\n", err)
		fmt.Println()
		fmt.Println("=== Workflow Failed ===")
		fmt.Printf("Error: %v\n", err)
		fmt.Println()
		fmt.Println("Check the Temporal UI for details:")
		fmt.Printf("  http://localhost:8233/namespaces/default/workflows/%s/%s\n",
			workflowRun.GetID(), workflowRun.GetRunID())
		return
	}

	fmt.Println()
	fmt.Println("=== Workflow Completed Successfully ===")
	fmt.Printf("Final Status: %s\n", result)
	fmt.Println()

	// Option 2: Query workflow state (non-blocking)
	// This would be used to check workflow state without waiting for completion
	/*
		var queryResult string
		queryValue, err := c.QueryWorkflow(
			context.Background(),
			workflowRun.GetID(),
			workflowRun.GetRunID(),
			"getOrderStatus", // Query handler name
		)
		if err != nil {
			log.Println("Query failed", err)
		} else {
			err = queryValue.Get(&queryResult)
			log.Printf("Current order status: %s", queryResult)
		}
	*/

	// Option 3: Signal workflow (to change behavior)
	// This would be used to send data to a running workflow
	/*
		err = c.SignalWorkflow(
			context.Background(),
			workflowRun.GetID(),
			workflowRun.GetRunID(),
			"updateOrderSignal", // Signal name
			updateData,          // Signal payload
		)
	*/

	// Option 4: Cancel workflow
	/*
		err = c.CancelWorkflow(
			context.Background(),
			workflowRun.GetID(),
			workflowRun.GetRunID(),
		)
	*/

	fmt.Println("View complete workflow history in Temporal UI:")
	fmt.Printf("  http://localhost:8233/namespaces/default/workflows/%s/%s\n",
		workflowRun.GetID(), workflowRun.GetRunID())
}

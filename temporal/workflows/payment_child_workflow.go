package workflows

import (
	"fmt"
	"temporal/models"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

/*
CHILD WORKFLOW CONCEPTS:

1. Why use Child Workflows?
   - Encapsulate complex sub-processes with their own lifecycle
   - Enable independent retry policies from parent workflow
   - Allow sub-processes to be executed multiple times or in parallel
   - Provide better organization and reusability
   - Can be versioned independently

2. Child Workflow vs Activity:
   - Child Workflow: Complex business logic with multiple steps, needs orchestration
   - Activity: Single unit of work, typically I/O operation

3. When does a Child Workflow retry?
   - When the child workflow itself fails (workflow-level error)
   - Based on the ParentClosePolicy when parent fails
   - Independent of activity retries within the child workflow
*/

// PaymentWorkflow is a child workflow that handles the payment process
// It demonstrates:
// - Child workflow pattern
// - Two-phase commit (authorize + capture)
// - Compensation logic (SAGA pattern)
// - Workflow-level retry policy
func PaymentWorkflow(ctx workflow.Context, order models.Order) (*models.PaymentInfo, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Payment child workflow started", "OrderID", order.OrderID)

	// Create payment info
	paymentInfo := &models.PaymentInfo{
		PaymentID:     fmt.Sprintf("pay-%s", order.OrderID),
		OrderID:       order.OrderID,
		Amount:        order.TotalAmount,
		PaymentMethod: "credit_card",
		Status:        "pending",
	}

	// Configure activity options for payment activities
	// These timeouts apply to activities within this child workflow
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second, // Max time for activity to complete
		// ScheduleToClose includes time in queue + execution
		ScheduleToCloseTimeout: 1 * time.Minute,

		// Activity-level retry policy
		// Retries happen BEFORE the activity is considered failed
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    1 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    5, // Retry up to 5 times

			// Don't retry on specific errors
			NonRetryableErrorTypes: []string{
				"INSUFFICIENT_FUNDS", // Business logic error
			},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Step 1: Authorize Payment
	// This holds the funds but doesn't charge yet
	var authErr error
	err := workflow.ExecuteActivity(ctx, "AuthorizePayment", paymentInfo).Get(ctx, &authErr)
	if err != nil {
		logger.Error("Payment authorization failed", "Error", err)
		// No compensation needed as we haven't done anything yet
		return nil, fmt.Errorf("payment authorization failed: %w", err)
	}

	logger.Info("Payment authorized", "AuthorizationID", paymentInfo.AuthorizationID)

	// Step 2: Capture Payment
	// This actually charges the customer
	var captureErr error
	err = workflow.ExecuteActivity(ctx, "CapturePayment", paymentInfo).Get(ctx, &captureErr)
	if err != nil {
		logger.Error("Payment capture failed", "Error", err)

		// COMPENSATION: Void the authorization since capture failed
		logger.Info("Initiating payment compensation (void authorization)")

		// Use a separate context with longer timeout for compensation
		compensationCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Second,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval: 1 * time.Second,
				MaximumAttempts: 3,
			},
		})

		var voidErr error
		compensationErr := workflow.ExecuteActivity(compensationCtx, "VoidAuthorization", paymentInfo).Get(compensationCtx, &voidErr)
		if compensationErr != nil {
			logger.Error("Compensation failed (void authorization)", "Error", compensationErr)
			// In production, this would trigger alerts
		}

		return nil, fmt.Errorf("payment capture failed: %w", err)
	}

	logger.Info("Payment workflow completed successfully",
		"PaymentID", paymentInfo.PaymentID,
		"Status", paymentInfo.Status)

	return paymentInfo, nil
}

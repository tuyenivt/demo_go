package workflows

import (
	"fmt"
	"temporal/models"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	enumspb "go.temporal.io/api/enums/v1"
)

/*
WORKFLOW DETERMINISM:

Workflows in Temporal MUST be deterministic because they can be replayed.
When a workflow is loaded from history, Temporal replays all the decisions.

DO in workflows:
- Call activities
- Call child workflows
- Use workflow.Now() for time
- Use workflow.Sleep() for delays
- Use workflow.GetLogger()
- Use workflow.SideEffect() for non-deterministic operations

DON'T in workflows:
- Make HTTP calls directly
- Access databases directly
- Use time.Now() or time.Sleep()
- Use random number generators
- Read from disk
- Any I/O operations

RETRY POLICIES EXPLAINED:

1. Activity Retry Policy:
   - Applied when an activity fails
   - Retries happen automatically before the workflow continues
   - Configured per activity or per workflow (default for all activities)
   - Example: Network timeout, temporary service unavailability

2. Workflow Retry Policy:
   - Applied when the ENTIRE workflow fails
   - Starts the workflow from the beginning
   - Configured when starting the workflow
   - Example: Business logic error that requires full restart
*/

// OrderWorkflow is the main workflow that orchestrates order processing
// This demonstrates the complete order fulfillment process with proper error handling
func OrderWorkflow(ctx workflow.Context, order models.Order) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Order workflow started", "OrderID", order.OrderID)

	// Update order status
	order.Status = "processing"

	// Track what we've done for compensation
	var inventoryReservation *models.InventoryReservation
	var paymentInfo *models.PaymentInfo

	// Configure default activity options
	// These apply to all activities unless overridden
	activityOptions := workflow.ActivityOptions{
		// StartToCloseTimeout: Maximum time for activity execution
		// This is the most important timeout for most activities
		StartToCloseTimeout: 30 * time.Second,

		// ScheduleToCloseTimeout: Maximum time from scheduling to completion
		// Includes time waiting in queue + execution time
		ScheduleToCloseTimeout: 1 * time.Minute,

		// ACTIVITY RETRY POLICY
		// This defines how activities are retried on failure
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    1 * time.Second,  // First retry after 1 second
			BackoffCoefficient: 2.0,              // Double the interval each retry
			MaximumInterval:    30 * time.Second, // Cap the retry interval
			MaximumAttempts:    5,                // Try up to 5 times total
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// STEP 1: Validate Order
	// Activity retries will happen automatically if this fails with transient errors
	err := workflow.ExecuteActivity(ctx, "ValidateOrder", order).Get(ctx, nil)
	if err != nil {
		logger.Error("Order validation failed", "Error", err)
		order.Status = "validation_failed"
		return order.Status, fmt.Errorf("order validation failed: %w", err)
	}
	logger.Info("Order validated successfully")

	// STEP 2: Process Payment using Child Workflow
	// Child workflows have their own lifecycle and retry policy
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowID: fmt.Sprintf("payment-%s", order.OrderID),

		// TaskQueue can be different from parent
		// TaskQueue: "payment-task-queue",

		// WORKFLOW RETRY POLICY
		// This is different from activity retry!
		// If the child workflow fails completely, it will be retried
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    2 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    1 * time.Minute,
			MaximumAttempts:    3, // Retry the entire workflow up to 3 times
		},

		// ParentClosePolicy: What happens to child when parent fails
		// - ABANDON: Child continues running
		// - REQUEST_CANCEL: Child receives cancellation request
		// - TERMINATE: Child is terminated immediately
		ParentClosePolicy: enumspb.PARENT_CLOSE_POLICY_TERMINATE,

		// Timeout for the entire child workflow
		WorkflowExecutionTimeout: 5 * time.Minute,
	}

	childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)

	// Execute child workflow
	var paymentResult *models.PaymentInfo
	childWorkflowFuture := workflow.ExecuteChildWorkflow(childCtx, PaymentWorkflow, order)

	err = childWorkflowFuture.Get(childCtx, &paymentResult)
	if err != nil {
		logger.Error("Payment workflow failed", "Error", err)
		order.Status = "payment_failed"
		// No compensation needed - child workflow handles its own compensation
		return order.Status, fmt.Errorf("payment processing failed: %w", err)
	}

	paymentInfo = paymentResult
	logger.Info("Payment processed successfully", "PaymentID", paymentInfo.PaymentID)

	// STEP 3: Reserve Inventory
	// If this fails, we need to compensate the payment
	err = workflow.ExecuteActivity(ctx, "ReserveInventory", order).Get(ctx, &inventoryReservation)
	if err != nil {
		logger.Error("Inventory reservation failed", "Error", err)
		order.Status = "inventory_failed"

		// COMPENSATION: Refund the payment
		logger.Info("Initiating compensation: refunding payment")
		compensationErr := compensatePayment(ctx, paymentInfo)
		if compensationErr != nil {
			logger.Error("Payment compensation failed", "Error", compensationErr)
			// In production, this would trigger manual intervention alerts
		}

		return order.Status, fmt.Errorf("inventory reservation failed: %w", err)
	}
	logger.Info("Inventory reserved successfully", "ReservationID", inventoryReservation.ReservationID)

	// STEP 4: Send Order Confirmation
	// This is a non-critical step - if it fails, we don't rollback everything
	notificationCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 20 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: 1 * time.Second,
			MaximumAttempts: 3, // Only retry 3 times for notifications
		},
	})

	notification := models.NotificationRequest{
		OrderID:     order.OrderID,
		CustomerID:  order.CustomerID,
		MessageType: "order_confirmation",
		Message:     fmt.Sprintf("Your order %s has been confirmed!", order.OrderID),
	}

	err = workflow.ExecuteActivity(notificationCtx, "SendNotification", notification).Get(notificationCtx, nil)
	if err != nil {
		// Log but don't fail the workflow
		logger.Warn("Failed to send notification (non-critical)", "Error", err)
		// In production, this could trigger a separate notification retry workflow
	} else {
		logger.Info("Notification sent successfully")
	}

	// WORKFLOW COMPLETED SUCCESSFULLY
	order.Status = "completed"
	logger.Info("Order workflow completed successfully",
		"OrderID", order.OrderID,
		"Status", order.Status,
		"PaymentID", paymentInfo.PaymentID,
		"ReservationID", inventoryReservation.ReservationID)

	return order.Status, nil
}

// compensatePayment is a helper function to refund payment
// This demonstrates the SAGA pattern for distributed transactions
func compensatePayment(ctx workflow.Context, paymentInfo *models.PaymentInfo) error {
	logger := workflow.GetLogger(ctx)

	if paymentInfo == nil {
		return nil
	}

	logger.Info("Compensating payment", "PaymentID", paymentInfo.PaymentID)

	compensationOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: 1 * time.Second,
			MaximumAttempts: 5, // Compensation is critical, retry more
		},
	}

	compensationCtx := workflow.WithActivityOptions(ctx, compensationOptions)

	return workflow.ExecuteActivity(compensationCtx, "RefundPayment", paymentInfo).Get(compensationCtx, nil)
}

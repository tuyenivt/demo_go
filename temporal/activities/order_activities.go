package activities

import (
	"context"
	"fmt"
	"math/rand"
	"temporal/models"
	"time"

	"go.temporal.io/sdk/activity"
)

// OrderActivities contains order-related activities
// Activities are where ALL side effects and I/O operations should happen
// They are non-deterministic and can be retried independently
type OrderActivities struct{}

// ValidateOrder checks if the order is valid
// This activity demonstrates:
// - Simple validation logic
// - Activity retry on transient failures
// - Activity heartbeats for long-running operations
func (a *OrderActivities) ValidateOrder(ctx context.Context, order models.Order) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Validating order", "OrderID", order.OrderID)

	// Simulate validation checks
	if order.CustomerID == "" {
		return fmt.Errorf("invalid order: customer ID is required")
	}

	if len(order.Items) == 0 {
		return fmt.Errorf("invalid order: no items in order")
	}

	if order.TotalAmount <= 0 {
		return fmt.Errorf("invalid order: invalid total amount")
	}

	// Simulate potential transient failure (network call to validation service)
	// This will be retried by the activity retry policy
	if rand.Float32() < 0.2 { // 20% failure rate for demo
		logger.Warn("Validation service temporarily unavailable")
		return fmt.Errorf("validation service unavailable: transient error")
	}

	// Simulate processing time
	time.Sleep(500 * time.Millisecond)

	logger.Info("Order validated successfully", "OrderID", order.OrderID)
	return nil
}

// ReserveInventory reserves inventory for the order
// This activity demonstrates:
// - Idempotent operations (can be safely retried)
// - Activity-level retry for transient failures
// - Using activity info for correlation
func (a *OrderActivities) ReserveInventory(ctx context.Context, order models.Order) (*models.InventoryReservation, error) {
	logger := activity.GetLogger(ctx)
	activityInfo := activity.GetInfo(ctx)

	logger.Info("Reserving inventory", "OrderID", order.OrderID, "AttemptCount", activityInfo.Attempt)

	// Simulate inventory check
	for _, item := range order.Items {
		logger.Info("Checking inventory", "ProductID", item.ProductID, "Quantity", item.Quantity)

		// Simulate out of stock scenario occasionally
		if item.ProductID == "prod-unavailable" {
			return nil, fmt.Errorf("product %s is out of stock", item.ProductID)
		}
	}

	// Simulate transient database failure
	if rand.Float32() < 0.15 { // 15% failure rate
		logger.Warn("Database temporarily unavailable during inventory reservation")
		return nil, fmt.Errorf("database error: connection timeout")
	}

	// Simulate processing time
	time.Sleep(300 * time.Millisecond)

	reservation := &models.InventoryReservation{
		ReservationID: fmt.Sprintf("res-%d", time.Now().Unix()),
		OrderID:       order.OrderID,
		Items:         order.Items,
		Status:        "reserved",
	}

	logger.Info("Inventory reserved successfully", "ReservationID", reservation.ReservationID)
	return reservation, nil
}

// SendNotification sends a notification to the customer
// This activity demonstrates:
// - Non-critical operations with retry
// - Different retry policies based on operation type
// - Logging for observability
func (a *OrderActivities) SendNotification(ctx context.Context, notif models.NotificationRequest) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending notification", "OrderID", notif.OrderID, "Type", notif.MessageType)

	// Simulate notification service call
	if rand.Float32() < 0.1 { // 10% failure rate
		logger.Warn("Notification service temporarily unavailable")
		return fmt.Errorf("notification service error: rate limit exceeded")
	}

	// Simulate sending email/SMS
	time.Sleep(200 * time.Millisecond)

	logger.Info("Notification sent successfully",
		"OrderID", notif.OrderID,
		"CustomerID", notif.CustomerID,
		"Type", notif.MessageType)

	return nil
}

// CompensateInventory releases reserved inventory (compensation logic)
// This is called when the workflow fails and we need to rollback
func (a *OrderActivities) CompensateInventory(ctx context.Context, reservation *models.InventoryReservation) error {
	logger := activity.GetLogger(ctx)

	if reservation == nil {
		logger.Info("No inventory reservation to compensate")
		return nil
	}

	logger.Info("Compensating inventory reservation", "ReservationID", reservation.ReservationID)

	// Simulate releasing inventory
	time.Sleep(100 * time.Millisecond)

	logger.Info("Inventory reservation released", "ReservationID", reservation.ReservationID)
	return nil
}

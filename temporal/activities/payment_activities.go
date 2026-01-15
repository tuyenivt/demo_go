package activities

import (
	"context"
	"fmt"
	"math/rand"
	"temporal/models"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

// PaymentActivities contains payment-related activities
type PaymentActivities struct{}

// AuthorizePayment authorizes the payment (holds the amount)
// This activity demonstrates:
// - Critical operation with strict retry policy
// - Idempotent operation design
// - Different failure modes
func (a *PaymentActivities) AuthorizePayment(ctx context.Context, paymentInfo *models.PaymentInfo) error {
	logger := activity.GetLogger(ctx)
	activityInfo := activity.GetInfo(ctx)

	logger.Info("Authorizing payment",
		"PaymentID", paymentInfo.PaymentID,
		"Amount", paymentInfo.Amount,
		"Attempt", activityInfo.Attempt)

	// Simulate different failure scenarios
	failureRate := rand.Float32()

	// Permanent failure: insufficient funds (should NOT retry)
	if paymentInfo.Amount > 10000 && failureRate < 0.05 {
		logger.Error("Payment authorization failed: insufficient funds")
		return temporal.NewApplicationError(
			"insufficient funds",
			"INSUFFICIENT_FUNDS", // Error type for conditional retry logic
			nil,
		)
	}

	// Transient failure: payment gateway timeout (SHOULD retry)
	if failureRate < 0.2 { // 20% transient failure rate
		logger.Warn("Payment gateway timeout")
		return fmt.Errorf("payment gateway timeout: transient error")
	}

	// Simulate payment gateway call
	time.Sleep(400 * time.Millisecond)

	// Update payment info with authorization ID
	paymentInfo.AuthorizationID = fmt.Sprintf("auth-%d", time.Now().Unix())
	paymentInfo.Status = "authorized"

	logger.Info("Payment authorized successfully",
		"AuthorizationID", paymentInfo.AuthorizationID)

	return nil
}

// CapturePayment captures the previously authorized payment
// This activity demonstrates:
// - Two-phase commit pattern
// - Dependency on previous activity result
// - Compensation-aware design
func (a *PaymentActivities) CapturePayment(ctx context.Context, paymentInfo *models.PaymentInfo) error {
	logger := activity.GetLogger(ctx)

	logger.Info("Capturing payment",
		"PaymentID", paymentInfo.PaymentID,
		"AuthorizationID", paymentInfo.AuthorizationID)

	// Verify authorization exists
	if paymentInfo.AuthorizationID == "" {
		return fmt.Errorf("cannot capture payment: no authorization ID")
	}

	// Simulate transient failure
	if rand.Float32() < 0.15 { // 15% failure rate
		logger.Warn("Payment capture temporarily failed")
		return fmt.Errorf("payment capture error: gateway unavailable")
	}

	// Simulate payment capture
	time.Sleep(300 * time.Millisecond)

	// Update payment info
	paymentInfo.CaptureID = fmt.Sprintf("cap-%d", time.Now().Unix())
	paymentInfo.Status = "captured"

	logger.Info("Payment captured successfully",
		"CaptureID", paymentInfo.CaptureID)

	return nil
}

// RefundPayment refunds a captured payment (compensation logic)
// This is called when we need to rollback a successful payment
func (a *PaymentActivities) RefundPayment(ctx context.Context, paymentInfo *models.PaymentInfo) error {
	logger := activity.GetLogger(ctx)

	if paymentInfo.CaptureID == "" {
		logger.Info("No captured payment to refund")
		return nil
	}

	logger.Info("Refunding payment",
		"PaymentID", paymentInfo.PaymentID,
		"CaptureID", paymentInfo.CaptureID)

	// Simulate refund processing
	time.Sleep(200 * time.Millisecond)

	paymentInfo.Status = "refunded"

	logger.Info("Payment refunded successfully")
	return nil
}

// VoidAuthorization voids an authorization (releases the hold)
// This is called when authorization succeeded but we didn't capture
func (a *PaymentActivities) VoidAuthorization(ctx context.Context, paymentInfo *models.PaymentInfo) error {
	logger := activity.GetLogger(ctx)

	if paymentInfo.AuthorizationID == "" {
		logger.Info("No authorization to void")
		return nil
	}

	logger.Info("Voiding authorization",
		"PaymentID", paymentInfo.PaymentID,
		"AuthorizationID", paymentInfo.AuthorizationID)

	// Simulate void processing
	time.Sleep(150 * time.Millisecond)

	paymentInfo.Status = "voided"

	logger.Info("Authorization voided successfully")
	return nil
}

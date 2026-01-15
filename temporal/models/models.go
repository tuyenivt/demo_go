package models

type Order struct {
	OrderID     string
	CustomerID  string
	Items       []OrderItem
	TotalAmount float64
	Status      string
}

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

type PaymentInfo struct {
	PaymentID       string
	OrderID         string
	Amount          float64
	PaymentMethod   string
	AuthorizationID string
	CaptureID       string
	Status          string
}

type InventoryReservation struct {
	ReservationID string
	OrderID       string
	Items         []OrderItem
	Status        string
}

type NotificationRequest struct {
	OrderID     string
	CustomerID  string
	MessageType string
	Message     string
}

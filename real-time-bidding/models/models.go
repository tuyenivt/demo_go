package models

import "time"

// BidRequest represents an incoming bid request
type BidRequest struct {
	RequestID    string            `json:"request_id"`
	UserID       string            `json:"user_id"`
	Timestamp    time.Time         `json:"timestamp"`
	Device       DeviceInfo        `json:"device"`
	UserSegments []string          `json:"user_segments"`
	Context      map[string]string `json:"context"`
}

// DeviceInfo contains device-specific information
type DeviceInfo struct {
	DeviceType string `json:"device_type"`
	OS         string `json:"os"`
	Browser    string `json:"browser"`
	IP         string `json:"ip"`
}

// Campaign represents an advertising campaign
type Campaign struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Budget      float64   `json:"budget"`
	BidPrice    float64   `json:"bid_price"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Targeting   Targeting `json:"targeting"`
	Status      string    `json:"status"`
	DailyBudget float64   `json:"daily_budget"`
}

// Targeting defines campaign targeting rules
type Targeting struct {
	UserSegments []string          `json:"user_segments"`
	DeviceTypes  []string          `json:"device_types"`
	OS           []string          `json:"os"`
	Browsers     []string          `json:"browsers"`
	CustomRules  map[string]string `json:"custom_rules"`
}

// BidResponse represents the response to a bid request
type BidResponse struct {
	RequestID  string  `json:"request_id"`
	CampaignID string  `json:"campaign_id"`
	Price      float64 `json:"price"`
	AdID       string  `json:"ad_id"`
	AdMarkup   string  `json:"ad_markup"`
	Latency    int64   `json:"latency_ms"`
}

// RedisKeyPatterns defines Redis key patterns
const (
	UserSegmentsKeyPattern   = "user:segments:%s"
	CampaignKeyPattern       = "campaign:%s"
	CampaignBudgetKeyPattern = "campaign:budget:%s"
	UserProfileKeyPattern    = "user:profile:%s"
)

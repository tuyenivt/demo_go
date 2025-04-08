package auction

import (
	"context"
	"testing"
	"time"

	"real-time-bidding/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedisCache is a mock implementation of the Redis cache
type MockRedisCache struct {
	mock.Mock
}

func (m *MockRedisCache) GetUserSegments(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisCache) GetCampaign(ctx context.Context, campaignID string) (*models.Campaign, error) {
	args := m.Called(ctx, campaignID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Campaign), args.Error(1)
}

func (m *MockRedisCache) GetCampaignBudget(ctx context.Context, campaignID string) (float64, error) {
	args := m.Called(ctx, campaignID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockRedisCache) UpdateCampaignBudget(ctx context.Context, campaignID string, amount float64) error {
	args := m.Called(ctx, campaignID, amount)
	return args.Error(0)
}

func TestProcessBidRequest(t *testing.T) {
	mockCache := new(MockRedisCache)
	service := NewAuctionService(mockCache)

	// Test data
	req := &models.BidRequest{
		RequestID: "test_req",
		UserID:    "test_user",
		Timestamp: time.Now(),
		Device: models.DeviceInfo{
			DeviceType: "mobile",
			OS:         "ios",
			Browser:    "safari",
			IP:         "127.0.0.1",
		},
		UserSegments: []string{"premium"},
	}

	campaign := &models.Campaign{
		ID:       "campaign1",
		Name:     "Test Campaign",
		Budget:   1000,
		BidPrice: 2.50,
		Status:   "active",
		Targeting: models.Targeting{
			DeviceTypes:  []string{"mobile"},
			OS:           []string{"ios"},
			UserSegments: []string{"premium"},
		},
	}

	// Set up expectations
	mockCache.On("GetUserSegments", mock.Anything, req.UserID).Return([]string{"premium"}, nil)
	mockCache.On("GetCampaign", mock.Anything, "campaign1").Return(campaign, nil)
	mockCache.On("GetCampaignBudget", mock.Anything, "campaign1").Return(100.0, nil)
	mockCache.On("UpdateCampaignBudget", mock.Anything, "campaign1", 2.50).Return(nil)

	// Process bid request
	ctx := context.Background()
	response, err := service.ProcessBidRequest(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "campaign1", response.CampaignID)
	assert.Equal(t, 2.50, response.Price)
	assert.Equal(t, "test_req", response.RequestID)

	mockCache.AssertExpectations(t)
}

func TestMatchesTargeting(t *testing.T) {
	service := NewAuctionService(nil)

	tests := []struct {
		name      string
		request   *models.BidRequest
		targeting *models.Targeting
		want      bool
	}{
		{
			name: "matches all criteria",
			request: &models.BidRequest{
				Device: models.DeviceInfo{
					DeviceType: "mobile",
					OS:         "ios",
				},
				UserSegments: []string{"premium"},
			},
			targeting: &models.Targeting{
				DeviceTypes:  []string{"mobile"},
				OS:           []string{"ios"},
				UserSegments: []string{"premium"},
			},
			want: true,
		},
		{
			name: "mismatched device type",
			request: &models.BidRequest{
				Device: models.DeviceInfo{
					DeviceType: "desktop",
					OS:         "ios",
				},
				UserSegments: []string{"premium"},
			},
			targeting: &models.Targeting{
				DeviceTypes:  []string{"mobile"},
				OS:           []string{"ios"},
				UserSegments: []string{"premium"},
			},
			want: false,
		},
		{
			name: "mismatched user segment",
			request: &models.BidRequest{
				Device: models.DeviceInfo{
					DeviceType: "mobile",
					OS:         "ios",
				},
				UserSegments: []string{"basic"},
			},
			targeting: &models.Targeting{
				DeviceTypes:  []string{"mobile"},
				OS:           []string{"ios"},
				UserSegments: []string{"premium"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.matchesTargeting(tt.request, tt.targeting)
			assert.Equal(t, tt.want, got)
		})
	}
}

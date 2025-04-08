package cache

import (
	"context"
	"testing"
	"time"

	"real-time-bidding/models"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisCache(t *testing.T) {
	// Create Redis client for testing
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // Use a different DB for testing
	})

	// Clean up test data
	ctx := context.Background()
	client.FlushDB(ctx)

	// Create cache instance
	cache, err := NewRedisCache("localhost:6379", "", 1)
	assert.NoError(t, err)
	defer cache.Close()

	// Test caching user segments
	t.Run("CacheUserSegments", func(t *testing.T) {
		userID := "test_user"
		segments := []string{"premium", "gaming"}

		err := cache.CacheUserSegments(ctx, userID, segments)
		assert.NoError(t, err)

		// Verify segments were cached
		cachedSegments, err := cache.GetUserSegments(ctx, userID)
		assert.NoError(t, err)
		assert.ElementsMatch(t, segments, cachedSegments)

		// Verify TTL
		ttl, err := client.TTL(ctx, "user:segments:"+userID).Result()
		assert.NoError(t, err)
		assert.True(t, ttl > 0 && ttl <= 5*time.Minute)
	})

	// Test caching campaign
	t.Run("CacheCampaign", func(t *testing.T) {
		campaign := &models.Campaign{
			ID:       "test_campaign",
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

		err := cache.CacheCampaign(ctx, campaign)
		assert.NoError(t, err)

		// Verify campaign was cached
		cachedCampaign, err := cache.GetCampaign(ctx, campaign.ID)
		assert.NoError(t, err)
		assert.Equal(t, campaign.ID, cachedCampaign.ID)
		assert.Equal(t, campaign.Name, cachedCampaign.Name)
		assert.Equal(t, campaign.Budget, cachedCampaign.Budget)

		// Verify budget was cached
		budget, err := cache.GetCampaignBudget(ctx, campaign.ID)
		assert.NoError(t, err)
		assert.Equal(t, campaign.DailyBudget, budget)

		// Verify TTL
		ttl, err := client.TTL(ctx, "campaign:"+campaign.ID).Result()
		assert.NoError(t, err)
		assert.True(t, ttl > 0 && ttl <= 5*time.Minute)
	})

	// Test budget updates
	t.Run("UpdateCampaignBudget", func(t *testing.T) {
		campaignID := "test_campaign"
		initialBudget := 100.0

		// Set initial budget
		err := client.Set(ctx, "campaign:budget:"+campaignID, initialBudget, 24*time.Hour).Err()
		assert.NoError(t, err)

		// Update budget
		err = cache.UpdateCampaignBudget(ctx, campaignID, 10.0)
		assert.NoError(t, err)

		// Verify updated budget
		budget, err := cache.GetCampaignBudget(ctx, campaignID)
		assert.NoError(t, err)
		assert.Equal(t, 90.0, budget)
	})
}

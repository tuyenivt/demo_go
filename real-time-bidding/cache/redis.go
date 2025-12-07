package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"real-time-bidding/models"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements the caching layer using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(addr string, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 100, // Adjust based on your needs
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// GetUserSegments retrieves user segments from Redis
func (c *RedisCache) GetUserSegments(ctx context.Context, userID string) ([]string, error) {
	key := fmt.Sprintf(models.UserSegmentsKeyPattern, userID)
	return c.client.SMembers(ctx, key).Result()
}

// GetCampaign retrieves a campaign from Redis
func (c *RedisCache) GetCampaign(ctx context.Context, campaignID string) (*models.Campaign, error) {
	key := fmt.Sprintf(models.CampaignKeyPattern, campaignID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var campaign models.Campaign
	if err := json.Unmarshal(data, &campaign); err != nil {
		return nil, err
	}

	return &campaign, nil
}

// GetCampaignBudget retrieves the current budget for a campaign
func (c *RedisCache) GetCampaignBudget(ctx context.Context, campaignID string) (float64, error) {
	key := fmt.Sprintf(models.CampaignBudgetKeyPattern, campaignID)
	return c.client.Get(ctx, key).Float64()
}

// UpdateCampaignBudget updates the campaign budget atomically
func (c *RedisCache) UpdateCampaignBudget(ctx context.Context, campaignID string, amount float64) error {
	key := fmt.Sprintf(models.CampaignBudgetKeyPattern, campaignID)
	pipe := c.client.Pipeline()

	// Get current budget
	currentBudget, err := pipe.Get(ctx, key).Float64()
	if err != nil {
		return err
	}

	// Update budget
	pipe.Set(ctx, key, currentBudget-amount, 24*time.Hour)

	_, err = pipe.Exec(ctx)
	return err
}

// CacheUserSegments caches user segments with TTL
func (c *RedisCache) CacheUserSegments(ctx context.Context, userID string, segments []string) error {
	key := fmt.Sprintf(models.UserSegmentsKeyPattern, userID)
	pipe := c.client.Pipeline()

	// Add segments to set
	pipe.SAdd(ctx, key, segments)

	// Set TTL (5 minutes)
	pipe.Expire(ctx, key, 5*time.Minute)

	_, err := pipe.Exec(ctx)
	return err
}

// CacheCampaign caches a campaign with TTL
func (c *RedisCache) CacheCampaign(ctx context.Context, campaign *models.Campaign) error {
	key := fmt.Sprintf(models.CampaignKeyPattern, campaign.ID)
	data, err := json.Marshal(campaign)
	if err != nil {
		return err
	}

	pipe := c.client.Pipeline()
	pipe.Set(ctx, key, data, 5*time.Minute)
	pipe.Set(ctx, fmt.Sprintf(models.CampaignBudgetKeyPattern, campaign.ID), campaign.DailyBudget, 24*time.Hour)

	_, err = pipe.Exec(ctx)
	return err
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

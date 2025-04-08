package redis

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"

	"github.com/go-redis/redis/v8"
)

// ShardedRedisClient represents a sharded Redis client
type ShardedRedisClient struct {
	clients []*redis.Client
	mu      sync.RWMutex
}

// NewShardedRedisClient creates a new sharded Redis client
func NewShardedRedisClient(hosts []string) (*ShardedRedisClient, error) {
	if len(hosts) == 0 {
		return nil, fmt.Errorf("no Redis hosts provided")
	}

	clients := make([]*redis.Client, len(hosts))
	for i, host := range hosts {
		client := redis.NewClient(&redis.Options{
			Addr:     host,
			Password: "", // Add password if needed
			DB:       0,
			PoolSize: 100, // Adjust pool size based on your needs
		})

		// Test connection
		ctx := context.Background()
		if err := client.Ping(ctx).Err(); err != nil {
			return nil, fmt.Errorf("failed to connect to Redis at %s: %w", host, err)
		}

		clients[i] = client
	}

	return &ShardedRedisClient{
		clients: clients,
	}, nil
}

// getShardIndex determines which Redis shard to use for a given key
func (c *ShardedRedisClient) getShardIndex(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) % len(c.clients)
}

// getClient returns the Redis client for a given key
func (c *ShardedRedisClient) getClient(key string) *redis.Client {
	shardIndex := c.getShardIndex(key)
	return c.clients[shardIndex]
}

// AddScore adds or updates a player's score
func (c *ShardedRedisClient) AddScore(ctx context.Context, playerID string, score float64) error {
	client := c.getClient(playerID)
	fmt.Sprintf("leaderboard:%s", playerID)
	return client.ZAdd(ctx, "leaderboard", &redis.Z{
		Score:  score,
		Member: playerID,
	}).Err()
}

// GetRank returns a player's rank (1-based)
func (c *ShardedRedisClient) GetRank(ctx context.Context, playerID string) (int64, error) {
	client := c.getClient(playerID)
	rank, err := client.ZRevRank(ctx, "leaderboard", playerID).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("player not found")
	}
	if err != nil {
		return 0, err
	}
	return rank + 1, nil
}

// GetTopN returns the top N players
func (c *ShardedRedisClient) GetTopN(ctx context.Context, n int64) ([]redis.Z, error) {
	// Since we need to get data from all shards, we'll use the first client
	// This is a limitation of the current implementation
	// For a production system, you might want to implement a more sophisticated
	// approach to get top N across shards
	client := c.clients[0]
	return client.ZRevRangeWithScores(ctx, "leaderboard", 0, n-1).Result()
}

// GetRange returns players within a rank range
func (c *ShardedRedisClient) GetRange(ctx context.Context, start, stop int64) ([]redis.Z, error) {
	client := c.clients[0]
	return client.ZRevRangeWithScores(ctx, "leaderboard", start, stop).Result()
}

// Close closes all Redis connections
func (c *ShardedRedisClient) Close() error {
	for _, client := range c.clients {
		if err := client.Close(); err != nil {
			return err
		}
	}
	return nil
}

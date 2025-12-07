package service

import (
	"context"
	"fmt"

	"real-time-gaming-leaderboard/internal/redis"
)

// LeaderboardService handles leaderboard operations
type LeaderboardService struct {
	redisClient *redis.ShardedRedisClient
}

// NewLeaderboardService creates a new leaderboard service
func NewLeaderboardService(redisClient *redis.ShardedRedisClient) *LeaderboardService {
	return &LeaderboardService{
		redisClient: redisClient,
	}
}

// Score represents a player's score
type Score struct {
	PlayerID string  `json:"player_id"`
	Score    float64 `json:"score"`
}

// PlayerRank represents a player's rank information
type PlayerRank struct {
	PlayerID string  `json:"player_id"`
	Score    float64 `json:"score"`
	Rank     int64   `json:"rank"`
}

// AddScore adds or updates a player's score
func (s *LeaderboardService) AddScore(ctx context.Context, playerID string, score float64) error {
	if playerID == "" {
		return fmt.Errorf("player_id is required")
	}

	return s.redisClient.AddScore(ctx, playerID, score)
}

// GetRank returns a player's rank
func (s *LeaderboardService) GetRank(ctx context.Context, playerID string) (*PlayerRank, error) {
	if playerID == "" {
		return nil, fmt.Errorf("player_id is required")
	}

	rank, err := s.redisClient.GetRank(ctx, playerID)
	if err != nil {
		return nil, err
	}

	// Get the player's score
	topPlayers, err := s.redisClient.GetTopN(ctx, rank)
	if err != nil {
		return nil, err
	}

	for _, player := range topPlayers {
		if player.Member == playerID {
			return &PlayerRank{
				PlayerID: playerID,
				Score:    player.Score,
				Rank:     rank,
			}, nil
		}
	}

	return nil, fmt.Errorf("player not found")
}

// GetTopN returns the top N players
func (s *LeaderboardService) GetTopN(ctx context.Context, n int64) ([]PlayerRank, error) {
	if n <= 0 {
		return nil, fmt.Errorf("n must be greater than 0")
	}

	topPlayers, err := s.redisClient.GetTopN(ctx, n)
	if err != nil {
		return nil, err
	}

	players := make([]PlayerRank, len(topPlayers))
	for i, player := range topPlayers {
		players[i] = PlayerRank{
			PlayerID: player.Member.(string),
			Score:    player.Score,
			Rank:     int64(i + 1),
		}
	}

	return players, nil
}

// GetRange returns players within a rank range
func (s *LeaderboardService) GetRange(ctx context.Context, start, stop int64) ([]PlayerRank, error) {
	if start < 0 || stop < start {
		return nil, fmt.Errorf("invalid range")
	}

	players, err := s.redisClient.GetRange(ctx, start, stop)
	if err != nil {
		return nil, err
	}

	result := make([]PlayerRank, len(players))
	for i, player := range players {
		result[i] = PlayerRank{
			PlayerID: player.Member.(string),
			Score:    player.Score,
			Rank:     start + int64(i) + 1,
		}
	}

	return result, nil
}

// Close closes the leaderboard service
func (s *LeaderboardService) Close() error {
	return s.redisClient.Close()
}

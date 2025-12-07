package handlers

import (
	"net/http"
	"real-time-gaming-leaderboard/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LeaderboardHandler handles HTTP requests for the leaderboard
type LeaderboardHandler struct {
	leaderboardService *service.LeaderboardService
}

// NewLeaderboardHandler creates a new leaderboard handler
func NewLeaderboardHandler(leaderboardService *service.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{
		leaderboardService: leaderboardService,
	}
}

// AddScore handles adding or updating a player's score
func (h *LeaderboardHandler) AddScore(c *gin.Context) {
	var score service.Score
	if err := c.ShouldBindJSON(&score); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.leaderboardService.AddScore(c.Request.Context(), score.PlayerID, score.Score); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "score updated successfully"})
}

// GetRank handles getting a player's rank
func (h *LeaderboardHandler) GetRank(c *gin.Context) {
	playerID := c.Param("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "player_id is required"})
		return
	}

	rank, err := h.leaderboardService.GetRank(c.Request.Context(), playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rank)
}

// GetTopN handles getting the top N players
func (h *LeaderboardHandler) GetTopN(c *gin.Context) {
	limitStr := c.Param("limit")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	players, err := h.leaderboardService.GetTopN(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, players)
}

// GetRange handles getting players within a rank range
func (h *LeaderboardHandler) GetRange(c *gin.Context) {
	startStr := c.Param("start")
	stopStr := c.Param("stop")

	start, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start parameter"})
		return
	}

	stop, err := strconv.ParseInt(stopStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stop parameter"})
		return
	}

	players, err := h.leaderboardService.GetRange(c.Request.Context(), start, stop)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, players)
}

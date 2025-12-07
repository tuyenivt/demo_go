package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"real-time-gaming-leaderboard/internal/config"
	"real-time-gaming-leaderboard/internal/handlers"
	"real-time-gaming-leaderboard/internal/middleware"
	"real-time-gaming-leaderboard/internal/redis"
	"real-time-gaming-leaderboard/internal/service"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Initialize Redis client
	redisClient, err := redis.NewShardedRedisClient(cfg.RedisHosts)
	if err != nil {
		logger.Fatal("Failed to create Redis client", zap.Error(err))
	}
	defer redisClient.Close()

	// Initialize services
	leaderboardService := service.NewLeaderboardService(redisClient)
	leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardService)

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit)

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(rateLimiter.RateLimit())

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/scores", leaderboardHandler.AddScore)
		v1.GET("/rank/:player_id", leaderboardHandler.GetRank)
		v1.GET("/top/:limit", leaderboardHandler.GetTopN)
		v1.GET("/range/:start/:stop", leaderboardHandler.GetRange)
	}

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.APIPort),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down server...")
		if err := srv.Close(); err != nil {
			logger.Error("Server forced to shutdown", zap.Error(err))
		}
	}()

	logger.Info("Starting server", zap.Int("port", cfg.APIPort))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}

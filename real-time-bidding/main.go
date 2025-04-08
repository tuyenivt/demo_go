package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"real-time-bidding/auction"
	"real-time-bidding/cache"
	"real-time-bidding/models"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
	}

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		0,
	)
	if err != nil {
		slog.Error("Failed to initialize Redis cache", "error", err)
		os.Exit(1)
	}
	defer redisCache.Close()

	// Initialize auction service
	auctionService := auction.NewAuctionService(redisCache)

	// Create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/bid", handleBidRequest(auctionService))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		slog.Info("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited properly")
}

func handleBidRequest(auctionService *auction.AuctionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var bidRequest models.BidRequest
		if err := json.NewDecoder(r.Body).Decode(&bidRequest); err != nil {
			slog.Error("Failed to decode bid request", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Process bid request
		ctx := r.Context()
		bidResponse, err := auctionService.ProcessBidRequest(ctx, &bidRequest)
		if err != nil {
			slog.Error("Failed to process bid request", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		if bidResponse == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if err := json.NewEncoder(w).Encode(bidResponse); err != nil {
			slog.Error("Failed to encode bid response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

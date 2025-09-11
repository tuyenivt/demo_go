package main

import (
	"basic-fullstack/internal/logger"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	logger := initLogger()

	http.HandleFunc("/health", healthCheck)
	http.Handle("/", http.FileServer(http.Dir("public")))

	server := &http.Server{
		Addr:         ":8080",
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("Server starting...")
	err := server.ListenAndServe()
	if err != nil {
		logger.Error("Server can't start. Error: %v", err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}

func initLogger() *logger.Logger {
	logger, err := logger.NewLogger("movie.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger. Error: %v", err)
	}
	defer logger.Close()
	return logger
}

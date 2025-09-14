package main

import (
	"basic-fullstack/internal/handlers"
	"basic-fullstack/internal/logger"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file was available")
	}
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		log.Fatal("DATABASE_URL not set")
	}
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Can't connect to database, error: %v", err)
	}
	defer db.Close()

	logger := initLogger()

	movieHandler := handlers.NewMovieHandler()

	http.HandleFunc("/api/movies/top", movieHandler.GetTopMovies)
	http.HandleFunc("/health", healthCheck)
	http.Handle("/", http.FileServer(http.Dir("public")))

	server := &http.Server{
		Addr:         ":8080",
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("Server starting...")
	err = server.ListenAndServe()
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

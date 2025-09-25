package main

import (
	"basic-fullstack/internal/data"
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

	http.HandleFunc("/health", healthCheck)

	accountRepo, err := data.NewAccountRepository(db, logger)
	if err != nil {
		log.Fatalf("Failed to initialize account repository: %v", err)
	}
	accountHandler := handlers.NewAccountHandler(accountRepo, logger)

	http.HandleFunc("/api/account/register/", accountHandler.Register)
	http.HandleFunc("/api/account/authenticate/", accountHandler.Authenticate)

	http.Handle("/api/account/favorites/", accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.GetFavorites)))
	http.Handle("/api/account/watchlist/", accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.GetWatchlist)))
	http.Handle("/api/account/save-to-collection/", accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.SaveToCollection)))

	movieRepository, err := data.NewMovieRepository(db, logger)
	if err != nil {
		log.Fatalf("Failed to initialize movie repository, error: %v", err)
	}
	movieHandler := handlers.NewMovieHandler(movieRepository, logger)

	http.HandleFunc("/api/movies/top", movieHandler.GetTopMovies)
	http.HandleFunc("/api/movies/search", movieHandler.SearchMoviesByName)
	http.HandleFunc("/api/movies/", movieHandler.GetMovieByID)
	http.HandleFunc("/api/genres", movieHandler.GetAllGenres)

	catchAllHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	}
	http.HandleFunc("/movies", catchAllHandler)
	http.HandleFunc("/movies/", catchAllHandler)
	http.HandleFunc("/account/", catchAllHandler)

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

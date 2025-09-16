package handlers

import (
	"basic-fullstack/internal/data"
	"basic-fullstack/internal/logger"
	"basic-fullstack/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

type MovieHandler struct {
	store  data.MovieStore
	logger *logger.Logger
}

func NewMovieHandler(store data.MovieStore, logger *logger.Logger) *MovieHandler {
	return &MovieHandler{store: store, logger: logger}
}

func (mh *MovieHandler) writeJSONResponse(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		mh.logger.Error("Failed to encode response, error: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (mh *MovieHandler) handleError(w http.ResponseWriter, err error, message string) bool {
	if err != nil {
		if err == data.ErrMovieNotFound {
			http.Error(w, message, http.StatusNotFound)
			return true
		}
		mh.logger.Error(message, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}
	return false
}

func (h *MovieHandler) parseID(w http.ResponseWriter, idStr string) (int, bool) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Invalid ID format", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func (mh *MovieHandler) GetTopMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := mh.store.GetTopMovies()
	if mh.handleError(w, err, "Failed to get top movies") {
		return
	}
	if mh.writeJSONResponse(w, movies) != nil {
		mh.logger.Info("Successfully served top movies")
	}
}

func (mh *MovieHandler) GetMovieByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/movies/"):]
	id, ok := mh.parseID(w, idStr)
	if !ok {
		return
	}

	movie, err := mh.store.GetMovieByID(id)
	if mh.handleError(w, err, "Failed to get movie by ID") {
		return
	}
	if mh.writeJSONResponse(w, movie) == nil {
		mh.logger.Info("Successfully served movie with ID: " + idStr)
	}
}

func (mh *MovieHandler) SearchMoviesByName(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	order := r.URL.Query().Get("order")
	genreStr := r.URL.Query().Get("genre")

	var genre *int
	if genreStr != "" {
		genreInt, ok := mh.parseID(w, genreStr)
		if !ok {
			return
		}
		genre = &genreInt
	}

	var movies []models.Movie
	var err error
	if query != "" {
		movies, err = mh.store.SearchMoviesByName(query, order, genre)
	}
	if mh.handleError(w, err, "Failed to search movies") {
		return
	}
	if mh.writeJSONResponse(w, movies) == nil {
		mh.logger.Info("Successfully served movies")
	}
}

func (mh *MovieHandler) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := mh.store.GetAllGenres()
	if mh.handleError(w, err, "Failed to get all genres") {
		return
	}
	if mh.writeJSONResponse(w, genres) == nil {
		mh.logger.Info("Successfully served genres")
	}
}

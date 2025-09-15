package handlers

import (
	"basic-fullstack/internal/data"
	"basic-fullstack/internal/logger"
	"encoding/json"
	"net/http"
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

func (mh *MovieHandler) GetTopMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := mh.store.GetTopMovies()
	if mh.handleError(w, err, "Failed to get top movies") {
		return
	}
	if mh.writeJSONResponse(w, movies) != nil {
		mh.logger.Info("Successfully served top movies")
	}
}

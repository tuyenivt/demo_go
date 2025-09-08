package api

import (
	"basic/internal/middleware"
	"basic/internal/store"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{workoutStore: workoutStore, logger: log.New(os.Stdout, "", log.Ldate|log.Ltime)}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.Error(w, "'id' param is required", http.StatusNotFound)
		return
	}

	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.Error(w, "'id' param must be integer", http.StatusBadRequest)
		return
	}
	wh.logger.Printf("Requesting workout data with id: %d\n", workoutID)

	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		http.Error(w, "Can't fetch workout data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		http.Error(w, "Invalid workout data", http.StatusBadRequest)
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser.IsAnonymous() {
		http.Error(w, "You must be logged in to create a workout", http.StatusUnauthorized)
		return
	}
	workout.UserID = currentUser.ID

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		http.Error(w, "Create workout failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdWorkout)
}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.Error(w, "'id' param is required", http.StatusNotFound)
		return
	}
	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.Error(w, "'id' param must be integer", http.StatusBadRequest)
		return
	}
	wh.logger.Printf("Requesting update workout data with id: %d\n", workoutID)

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		http.Error(w, "Can't fetch workout data", http.StatusInternalServerError)
		return
	}

	if existingWorkout == nil {
		http.Error(w, "Can't found workout data", http.StatusNotFound)
		return
	}

	var workout store.Workout
	err = json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		http.Error(w, "Invalid workout data", http.StatusBadRequest)
		return
	}

	existingWorkout.Title = workout.Title
	existingWorkout.Description = workout.Description
	existingWorkout.DurationMinutes = workout.DurationMinutes
	existingWorkout.CaloriesBurned = workout.CaloriesBurned
	existingWorkout.Entries = workout.Entries

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser.IsAnonymous() {
		http.Error(w, "You must be logged in to update this workout", http.StatusUnauthorized)
		return
	}
	if existingWorkout.UserID != currentUser.ID {
		http.Error(w, "You are not authorized to update this workout", http.StatusForbidden)
		return
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		http.Error(w, "Update workout failed", http.StatusInternalServerError)
		return
	}

	updatedWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		http.Error(w, "Can't fetch workout data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedWorkout)
}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.Error(w, "'id' param is required", http.StatusNotFound)
		return
	}

	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.Error(w, "'id' param must be integer", http.StatusBadRequest)
		return
	}
	wh.logger.Printf("Requesting delete workout data with id: %d\n", workoutID)

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser.IsAnonymous() {
		http.Error(w, "You must be logged in to delete this workout", http.StatusUnauthorized)
		return
	}

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutID)
	if err == sql.ErrNoRows {
		http.Error(w, "Can't found workout data", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Can't fetch workout data", http.StatusInternalServerError)
		return
	}
	if workoutOwner != currentUser.ID {
		http.Error(w, "You are not authorized to delete this workout", http.StatusForbidden)
		return
	}

	err = wh.workoutStore.DeleteWorkoutByID(workoutID)
	if err == sql.ErrNoRows {
		http.Error(w, "Can't found workout data", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Can't delete workout data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

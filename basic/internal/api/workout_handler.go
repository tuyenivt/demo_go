package api

import (
	"basic/internal/store"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{workoutStore: workoutStore}
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
	fmt.Printf("Requesting workout data with id: %d\n", workoutID)

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
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		http.Error(w, "Create workout failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdWorkout)
}

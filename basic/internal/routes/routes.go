package routes

import (
	"basic/internal/app"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)

	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkout)
	r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)

	r.Get("/users/{username}", app.UserHandler.HandleGetUserByUsername)
	r.Post("/users", app.UserHandler.HandleRegisterUser)

	r.Post("/tokens", app.TokenHandler.HandleCreateNewToken)
	r.Delete("/tokens", app.TokenHandler.HandleDeleteAllTokensByUserID)

	return r
}

package app

import (
	"basic/internal/api"
	"basic/internal/middleware"
	"basic/internal/store"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	Middleware     *middleware.UserMiddleware
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	workoutHandler := api.NewWorkoutHandler(workoutStore)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	middleware := middleware.NewUserMiddleware(userStore)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		Middleware:     middleware,
		DB:             pgDB,
	}
	return app, nil
}

func (*Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}

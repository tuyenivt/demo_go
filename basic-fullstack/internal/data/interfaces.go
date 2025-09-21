package data

import "basic-fullstack/internal/models"

type MovieStore interface {
	GetTopMovies() ([]models.Movie, error)
	GetMovieByID(id int) (models.Movie, error)
	SearchMoviesByName(name string, order string, genre *int) ([]models.Movie, error)
	GetAllGenres() ([]models.Genre, error)
}

type AccountStorage interface {
	Authenticate(string, string) (bool, error)
	Register(string, string, string) (bool, error)
	GetAccountDetails(string) (models.User, error)
	SaveCollection(models.User, int, string) (bool, error)
}

package data

import "basic-fullstack/internal/models"

type MovieStore interface {
	GetTopMovies() ([]models.Movie, error)
	GetMovieByID(id int) (models.Movie, error)
	SearchMoviesByName(name string) ([]models.Movie, error)
}

type GenreStore interface {
	GetAllGenres() ([]models.Genre, error)
}

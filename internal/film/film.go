package film

import (
	"context"
	"films_library/internal/model"
)

type Usecase interface {
	GetFilms(ctx context.Context, filter model.FilmFilter) ([]model.Film, error)
	AddFilm(ctx context.Context, film model.AddFilmRequest) (uint64, error)
	UpdateFilm(ctx context.Context, film model.Film) (uint64, error)
	DeleteFilm(ctx context.Context, id uint64) (uint64, error)
	GetFilm(ctx context.Context, id uint64) (model.Film, error)
	SearchFilm(ctx context.Context, search string) ([]model.Film, error)
}

type Repository interface {
	GetFilms(ctx context.Context, filter model.FilmFilter) ([]model.Film, error)
	GetFilm(ctx context.Context, id uint64) (model.Film, error)
	AddFilm(ctx context.Context, film model.AddFilmRequest) (uint64, error)
	UpdateFilm(ctx context.Context, film model.Film) (uint64, error)
	DeleteFilm(ctx context.Context, id uint64) (uint64, error)
	SearchFilm(ctx context.Context, search string) ([]model.Film, error)
}

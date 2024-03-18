package usecase

import (
	"context"

	"films_library/internal/film"
	"films_library/internal/model"
	"films_library/pkg/logger"
)

type FilmUsecase struct {
	FilmRepository film.Repository
	logger         logger.Interface
}

func NewFilmUsecase(fr film.Repository, l logger.Interface) *FilmUsecase {
	return &FilmUsecase{fr, l}
}

func (fu *FilmUsecase) GetFilms(ctx context.Context, filter model.FilmFilter) ([]model.Film, error) {
	films, err := fu.FilmRepository.GetFilms(ctx, filter)
	if err != nil {
		return []model.Film{}, err
	}
	return films, nil
}

func (fu *FilmUsecase) AddFilm(ctx context.Context, film model.AddFilmRequest) (uint64, error) {
	id, err := fu.FilmRepository.AddFilm(ctx, film)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (fu *FilmUsecase) UpdateFilm(ctx context.Context, film model.Film) (uint64, error) {
	id, err := fu.FilmRepository.UpdateFilm(ctx, film)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (fu *FilmUsecase) DeleteFilm(ctx context.Context, id uint64) (uint64, error) {
	id, err := fu.FilmRepository.DeleteFilm(ctx, id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (fu *FilmUsecase) GetFilm(ctx context.Context, id uint64) (model.Film, error) {
	film, err := fu.FilmRepository.GetFilm(ctx, id)
	if err != nil {
		return model.Film{}, err
	}
	return film, nil
}

func (fu *FilmUsecase) SearchFilm(ctx context.Context, search string) ([]model.Film, error) {
	films, err := fu.FilmRepository.SearchFilm(ctx, search)
	if err != nil {
		return []model.Film{}, err
	}
	return films, nil
}

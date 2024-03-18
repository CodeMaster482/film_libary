package usecase

import (
	"context"
	"errors"
	"films_library/internal/model"
	"films_library/pkg/logger"
	"testing"

	mock_film "films_library/internal/film/mocks"

	"github.com/golang/mock/gomock"
)

func TestFilmUsecase_GetFilms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.NewMockInterface(ctrl)
	mockRepo := mock_film.NewMockRepository(ctrl)
	mockUsecase := NewFilmUsecase(mockRepo, logger)

	ctx := context.Background()

	testCases := []struct {
		name          string
		filter        model.FilmFilter
		expectedFilms []model.Film
		expectedError error
	}{
		{
			name:          "Valid case",
			filter:        model.FilmFilter{},
			expectedFilms: []model.Film{{ID: 1, Title: "Film 1"}, {ID: 2, Title: "Film 2"}},
			expectedError: nil,
		},
		{
			name:          "Error from repository",
			filter:        model.FilmFilter{},
			expectedFilms: nil,
			expectedError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().GetFilms(ctx, tc.filter).Return(tc.expectedFilms, tc.expectedError)

			films, err := mockUsecase.GetFilms(ctx, tc.filter)
			if err != tc.expectedError {
				t.Errorf("Unexpected error: expected %v, got %v", tc.expectedError, err)
			}

			if len(films) != len(tc.expectedFilms) {
				t.Errorf("Expected %d films, got %d", len(tc.expectedFilms), len(films))
			}

			for i := range tc.expectedFilms {
				if films[i] != tc.expectedFilms[i] {
					t.Errorf("Expected film %v, got %v", tc.expectedFilms[i], films[i])
				}
			}
		})
	}
}

func TestFilmUsecase_AddFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.NewMockInterface(ctrl)
	mockRepo := mock_film.NewMockRepository(ctrl)
	mockUsecase := NewFilmUsecase(mockRepo, logger)

	ctx := context.Background()

	testCases := []struct {
		name          string
		filmToAdd     model.AddFilmRequest
		expectedID    uint64
		expectedError error
	}{
		{
			name:          "Valid film",
			filmToAdd:     model.AddFilmRequest{Title: "Test Film"},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name:          "Error from repository",
			filmToAdd:     model.AddFilmRequest{Title: "Invalid Film"},
			expectedID:    0,
			expectedError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().AddFilm(ctx, tc.filmToAdd).Return(tc.expectedID, tc.expectedError)

			id, err := mockUsecase.AddFilm(ctx, tc.filmToAdd)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("Unexpected error: expected %v, got %v", tc.expectedError, err)
			}

			if id != tc.expectedID {
				t.Errorf("Expected ID %d, got %d", tc.expectedID, id)
			}
		})
	}
}

func TestFilmUsecase_UpdateFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.NewMockInterface(ctrl)
	mockRepo := mock_film.NewMockRepository(ctrl)
	mockUsecase := NewFilmUsecase(mockRepo, logger)

	ctx := context.Background()

	testCases := []struct {
		name          string
		filmToUpdate  model.Film
		expectedID    uint64
		expectedError error
	}{
		{
			name:          "Valid film",
			filmToUpdate:  model.Film{ID: 1, Title: "Updated Film"},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name:          "Error from repository",
			filmToUpdate:  model.Film{ID: 2, Title: "Invalid Film"},
			expectedID:    0,
			expectedError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().UpdateFilm(ctx, tc.filmToUpdate).Return(tc.expectedID, tc.expectedError)

			id, err := mockUsecase.UpdateFilm(ctx, tc.filmToUpdate)
			if err != tc.expectedError {
				t.Errorf("Unexpected error: expected %v, got %v", tc.expectedError, err)
			}

			if id != tc.expectedID {
				t.Errorf("Expected ID %d, got %d", tc.expectedID, id)
			}
		})
	}
}

func TestFilmUsecase_DeleteFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.NewMockInterface(ctrl)
	mockRepo := mock_film.NewMockRepository(ctrl)
	mockUsecase := NewFilmUsecase(mockRepo, logger)

	ctx := context.Background()

	testCases := []struct {
		name          string
		filmID        uint64
		expectedID    uint64
		expectedError error
	}{
		{
			name:          "Valid film ID",
			filmID:        1,
			expectedID:    1,
			expectedError: nil,
		},
		{
			name:          "Error from repository",
			filmID:        2,
			expectedID:    0,
			expectedError: errors.New("repository error"),
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().DeleteFilm(ctx, tc.filmID).Return(tc.expectedID, tc.expectedError)

			id, err := mockUsecase.DeleteFilm(ctx, tc.filmID)
			if err != tc.expectedError {
				t.Errorf("Unexpected error: expected %v, got %v", tc.expectedError, err)
			}

			if id != tc.expectedID {
				t.Errorf("Expected ID %d, got %d", tc.expectedID, id)
			}
		})
	}
}

func TestFilmUsecase_GetFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.NewMockInterface(ctrl)
	mockRepo := mock_film.NewMockRepository(ctrl)
	mockUsecase := NewFilmUsecase(mockRepo, logger)

	ctx := context.Background()

	testCases := []struct {
		name          string
		filmID        uint64
		expectedFilm  model.Film
		expectedError error
	}{
		{
			name:          "Valid film ID",
			filmID:        1,
			expectedFilm:  model.Film{ID: 1, Title: "Test Film"},
			expectedError: nil,
		},
		{
			name:          "Error from repository",
			filmID:        2,
			expectedFilm:  model.Film{},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().GetFilm(ctx, tc.filmID).Return(tc.expectedFilm, tc.expectedError)

			film, err := mockUsecase.GetFilm(ctx, tc.filmID)
			if err != tc.expectedError {
				t.Errorf("Unexpected error: expected %v, got %v", tc.expectedError, err)
			}

			if film != tc.expectedFilm {
				t.Errorf("Expected film %v, got %v", tc.expectedFilm, film)
			}
		})
	}
}

func TestFilmUsecase_SearchFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.NewMockInterface(ctrl)
	mockRepo := mock_film.NewMockRepository(ctrl)
	mockUsecase := NewFilmUsecase(mockRepo, logger)

	ctx := context.Background()

	testCases := []struct {
		name          string
		search        string
		expectedFilms []model.Film
		expectedError error
	}{
		{
			name:          "Valid search term",
			search:        "Action",
			expectedFilms: []model.Film{{ID: 1, Title: "Action Film 1"}, {ID: 2, Title: "Action Film 2"}},
			expectedError: nil,
		},
		{
			name:          "Error from repository",
			search:        "Comedy",
			expectedFilms: nil,
			expectedError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().SearchFilm(ctx, tc.search).Return(tc.expectedFilms, tc.expectedError)

			films, err := mockUsecase.SearchFilm(ctx, tc.search)
			if err != tc.expectedError {
				t.Errorf("Unexpected error: expected %v, got %v", tc.expectedError, err)
			}

			if len(films) != len(tc.expectedFilms) {
				t.Errorf("Expected %d films, got %d", len(tc.expectedFilms), len(films))
			}

			for i := range tc.expectedFilms {
				if films[i] != tc.expectedFilms[i] {
					t.Errorf("Expected film %v, got %v", tc.expectedFilms[i], films[i])
				}
			}
		})
	}
}

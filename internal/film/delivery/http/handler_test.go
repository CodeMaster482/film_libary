package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock_film "films_library/internal/film/mocks"
	"films_library/internal/model"
	"films_library/pkg/logger"

	"github.com/golang/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestFilmHandler_GetFilms(t *testing.T) {
	MockResponse := []model.Film{{ID: 1, Title: "Forest Gamp", Description: "...", Rating: 10}}
	tests := []struct {
		name          string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mock_film.MockUsecase)
		queryParams   map[string]string
	}{
		{
			name:         "Successful call to GetFilms with null query",
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":[{"film_id":1,"title":"Forest Gamp","description":"...","release_date":"0001-01-01T00:00:00Z","rating":10}]}`,
			mockUsecaseFn: func(mockUsecase *mock_film.MockUsecase) {
				mockUsecase.EXPECT().GetFilms(gomock.Any(), gomock.Any()).Return(MockResponse, nil)
			},
			queryParams: map[string]string{},
		},
		{
			name:         "Successful rating query",
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":[{"film_id":1,"title":"Forest Gamp","description":"...","release_date":"0001-01-01T00:00:00Z","rating":10}]}`,
			mockUsecaseFn: func(mockUsecase *mock_film.MockUsecase) {
				mockUsecase.EXPECT().GetFilms(gomock.Any(), gomock.Any()).Return(MockResponse, nil)
			},
			queryParams: map[string]string{"sort_by": "rating", "sort_order": "desc"},
		},
		{
			name:         "Successful rating query",
			expectedCode: http.StatusOK,
			expectedBody: `{"status":200,"body":[{"film_id":1,"title":"Forest Gamp","description":"...","release_date":"0001-01-01T00:00:00Z","rating":10}]}`,
			mockUsecaseFn: func(mockUsecase *mock_film.MockUsecase) {
				mockUsecase.EXPECT().GetFilms(gomock.Any(), gomock.Any()).Return(MockResponse, nil)
			},
			queryParams: map[string]string{"sort_by": "rating", "sort_order": "desc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := logger.NewMockInterface(ctrl)
			mockUsecase := mock_film.NewMockUsecase(ctrl)
			tt.mockUsecaseFn(mockUsecase)

			handler := FilmHandler{filmUsecase: mockUsecase, logger: logger}

			req := httptest.NewRequest("GET", "/film", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			recorder := httptest.NewRecorder()

			handler.GetFilms(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

func TestFilmHandler_AddFilm(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   string
		expectedCode  int
		expectedBody  string
		mockUsecaseFn func(*mock_film.MockUsecase)
	}{
		{
			name:         "Successful film addition",
			requestBody:  `{"title":"Forest Gump","description":"...","release_date":"2024-03-18","rating":9.5}`,
			expectedCode: http.StatusCreated,
			expectedBody: `{"status":201,"body":1}`,
			mockUsecaseFn: func(mockUsecase *mock_film.MockUsecase) {
				mockUsecase.EXPECT().AddFilm(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
			},
		},
		{
			name:          "Invalid request body",
			requestBody:   `{"title":"Forest Gump","description":"...","release_date":"2024-03-18"}`,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"status":400,"message":"Invalid request"}`,
			mockUsecaseFn: func(mockUsecase *mock_film.MockUsecase) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := logger.NewMockInterface(ctrl)
			mockUsecase := mock_film.NewMockUsecase(ctrl)
			if tt.mockUsecaseFn != nil {
				tt.mockUsecaseFn(mockUsecase)
			}

			handler := FilmHandler{filmUsecase: mockUsecase, logger: logger}

			req := httptest.NewRequest("POST", "/film", strings.NewReader(tt.requestBody))
			recorder := httptest.NewRecorder()

			handler.AddFilm(recorder, req)

			actual := strings.TrimSpace(recorder.Body.String())

			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedBody, actual)
		})
	}
}

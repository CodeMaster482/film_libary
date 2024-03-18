package http

import (
	"net/http"
	"strconv"

	"films_library/internal/film"
	"films_library/internal/model"
	"films_library/pkg/logger"
	"films_library/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/mailru/easyjson"
)

type FilmHandler struct {
	filmUsecase film.Usecase
	logger      logger.Interface
}

func NewFilmHandler(mux *http.ServeMux, fu film.Usecase, l logger.Interface) {
	r := &FilmHandler{fu, l}

	mux.HandleFunc("/film", r.GetFilms)
	mux.HandleFunc("/film/add", r.AddFilm)
	mux.HandleFunc("/film/update", r.UpdateFilm)
	mux.HandleFunc("/film/delete", r.DeleteFilm)
	mux.HandleFunc("/film/search", r.SearchFilm)
}

func (h *FilmHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	sortBy := queryParams.Get("sort_by")
	sortOrder := queryParams.Get("sort_order")

	if sortBy == "" {
		sortBy = "rating"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	filter := model.FilmFilter{
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	v := validator.New()
	if err := v.Struct(filter); err != nil {
		filter.SortBy = "rating"
		filter.SortOrder = "desc"
	}

	films, err := h.filmUsecase.GetFilms(r.Context(), filter)
	if err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", h.logger)
		return
	}

	response.SuccessResponse(w, http.StatusOK, films)
}

func (h *FilmHandler) AddFilm(w http.ResponseWriter, r *http.Request) {
	var film model.AddFilmRequest
	if err := easyjson.UnmarshalFromReader(r.Body, &film); err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusBadRequest, "Corrupted request body", h.logger)
		return
	}

	v := validator.New()
	if err := v.Struct(film); err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", h.logger)
		return
	}

	id, err := h.filmUsecase.AddFilm(r.Context(), film)
	if err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", h.logger)
		return
	}
	response.SuccessResponse(w, http.StatusCreated, id)
}

func (h *FilmHandler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	var film model.Film
	if err := easyjson.UnmarshalFromReader(r.Body, &film); err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusBadRequest, "Corrupted request body", h.logger)
		return
	}

	v := validator.New()
	if err := v.Struct(film); err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", h.logger)
		return
	}

	id, err := h.filmUsecase.UpdateFilm(r.Context(), film)
	if err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", h.logger)
		return
	}
	response.SuccessResponse(w, http.StatusOK, id)
}

func (h *FilmHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	filmId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusBadRequest, "Bad query param", h.logger)
		return
	}

	if filmId < 0 {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", h.logger)
		return
	}

	id, err := h.filmUsecase.DeleteFilm(r.Context(), uint64(filmId))
	if err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", h.logger)
		return
	}
	response.SuccessResponse(w, http.StatusOK, id)
}

func (h *FilmHandler) SearchFilm(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	search := queryParams.Get("search")

	if search == "" {
		h.logger.Error("Empty title")
		response.ErrorResponse(w, http.StatusBadRequest, "Empty title", h.logger)
		return
	}

	film, err := h.filmUsecase.SearchFilm(r.Context(), search)
	if err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", h.logger)
		return
	}
	response.SuccessResponse(w, http.StatusOK, film)
}

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

// GetFilms handles the HTTP GET request to retrieve a list of films.
// It allows optional query parameters for sorting the results.
// @Summary Get films
// @Description Retrieves a list of films with optional sorting.
// @Tags films
// @Produce json
// @Param sort_by query string false "Field to sort by (e.g., 'rating')"
// @Param sort_order query string false "Sort order ('asc' for ascending or 'desc' for descending)"
// @Success 200 {array} model.Film "List of films"
// @Failure 500 {string} string "Internal Server Error"
// @Router /film [get]
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

// AddFilm handles the HTTP POST request to add a new film.
// @Summary Add film
// @Description Adds a new film to the system.
// @Tags films
// @Accept json
// @Produce json
// @Param film body model.AddFilmRequest true "Film object to be added"
// @Success 201 {string} string "ID of the newly added film"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /film/add [post]
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

// UpdateFilm handles the HTTP PUT request to update an existing film.
// @Summary Update film
// @Description Updates an existing film in the system.
// @Tags films
// @Accept json
// @Produce json
// @Param film body model.Film true "Film object to be updated"
// @Success 200 {string} string "ID of the updated film"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /film/update [put]
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

// DeleteFilm handles the HTTP DELETE request to delete an existing film by ID.
// @Summary Delete film
// @Description Deletes an existing film from the system by ID.
// @Tags films
// @Param id query integer true "ID of the film to be deleted"
// @Success 200 {string} string "ID of the deleted film"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /film/delete [delete]
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

// SearchFilm handles the HTTP GET request to search for films by title.
// @Summary Search film
// @Description Searches for films by title.
// @Tags films
// @Produce json
// @Param search query string true "Title to search for"
// @Success 200 {object} model.Film "Found film"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /film/search [get]
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

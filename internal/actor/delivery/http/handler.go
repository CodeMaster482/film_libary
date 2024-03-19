package http

import (
	"errors"
	"net/http"
	"strconv"

	"films_library/internal/actor"
	"films_library/internal/model"
	"films_library/pkg/logger"
	"films_library/pkg/response"

	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4"
	"github.com/mailru/easyjson"
)

type ActorHandler struct {
	actorUsecase actor.Usecase
	logger       logger.Interface
}

func NewActorHandler(mux *http.ServeMux, au actor.Usecase, l logger.Interface) {
	r := &ActorHandler{au, l}

	mux.HandleFunc("/actors", r.GetActor)
	mux.HandleFunc("/actors/add", r.AddActor)
	mux.HandleFunc("/actors/update", r.UpdateActor)
	mux.HandleFunc("/actors/delete", r.DeleteActor)
}

// GetActor handles the HTTP GET request to retrieve a list of actors.
// @Summary Get actors
// @Description Retrieves a list of actors.
// @Tags actors
// @Produce json
// @Success 200 {array} model.Actor "List of actors"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /actors [get]
func (h *ActorHandler) GetActor(w http.ResponseWriter, r *http.Request) {
	actors, err := h.actorUsecase.GetActors(r.Context())
	if err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", h.logger)
		return
	}

	response.SuccessResponse(w, http.StatusOK, actors)
}

// AddActor handles the HTTP POST request to add a new actor.
// @Summary Add actor
// @Description Adds a new actor to the system.
// @Tags actors
// @Accept json
// @Produce json
// @Param actor body model.Actor true "Actor object to be added"
// @Success 200 {string} string "ID of the newly added actor"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /actors/add [post]
func (h *ActorHandler) AddActor(w http.ResponseWriter, r *http.Request) {
	var actor model.Actor
	if err := easyjson.UnmarshalFromReader(r.Body, &actor); err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusBadRequest, "Corrupted request body", h.logger)
		return
	}

	v := validator.New()
	if err := v.Struct(actor); err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", h.logger)
		return
	}

	id, err := h.actorUsecase.AddActor(r.Context(), &actor)
	if err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", h.logger)
		return
	}

	response.SuccessResponse(w, http.StatusOK, id)
}

// UpdateActor handles the HTTP PUT request to update an existing actor.
// @Summary Update actor
// @Description Updates an existing actor in the system.
// @Tags actors
// @Accept json
// @Produce json
// @Param actor body model.Actor true "Actor object to be updated"
// @Success 200 {object} model.Actor "Updated actor object"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Object don't exist"
// @Failure 500 {string} string "Internal Server Error"
// @Router /actors/update [put]
func (h *ActorHandler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	var actor model.Actor
	if err := easyjson.UnmarshalFromReader(r.Body, &actor); err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusBadRequest, "Corrupted request body", h.logger)
		return
	}

	v := validator.New()
	if err := v.Struct(actor); err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", h.logger)
		return
	}

	updatedActor, err := h.actorUsecase.UpdateActor(r.Context(), &actor)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.logger.Info("user bad request: %s", err)
			response.ErrorResponse(w, http.StatusInternalServerError, "Object don't exist", h.logger)
			return
		}
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", h.logger)
		return
	}

	response.SuccessResponse(w, http.StatusOK, updatedActor)
}

// DeleteActor handles the HTTP DELETE request to delete an existing actor by ID.
// @Summary Delete actor
// @Description Deletes an existing actor from the system by ID.
// @Tags actors
// @Param id query integer true "ID of the actor to be deleted"
// @Success 200 {string} string "ID of the deleted actor"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Object don't exist"
// @Failure 500 {string} string "Internal Server Error"
// @Router /actors/delete [delete]
func (h *ActorHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	actorId, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusBadRequest, "Bad query param", h.logger)
		return
	}

	if actorId < 0 {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", h.logger)
		return
	}

	id, err := h.actorUsecase.DeleteActor(r.Context(), 1)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.logger.Info("user bad request: %s", err)
			response.ErrorResponse(w, http.StatusInternalServerError, "Object don't exist", h.logger)
			return
		}
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", h.logger)
		return
	}

	response.SuccessResponse(w, http.StatusOK, id)
}

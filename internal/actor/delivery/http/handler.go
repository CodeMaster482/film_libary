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

	mux.HandleFunc("/actor", r.GetActor)
	mux.HandleFunc("/actor/add", r.AddActor)
	mux.HandleFunc("/actor/update", r.UpdateActor)
	mux.HandleFunc("/actor/delete", r.DeleteActor)
}

func (h *ActorHandler) GetActor(w http.ResponseWriter, r *http.Request) {
	actors, err := h.actorUsecase.GetActors(r.Context())
	if err != nil {
		h.logger.Error(err)
		response.ErrorResponse(w, http.StatusInternalServerError, "Internal server error", h.logger)
		return
	}

	response.SuccessResponse(w, http.StatusOK, actors)
}

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

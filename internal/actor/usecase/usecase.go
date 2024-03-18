package usecase

import (
	"context"
	"films_library/internal/actor"
	"films_library/internal/model"
	"films_library/pkg/logger"
)

type Usecase struct {
	actorRepo actor.Repository
	logger    logger.Interface
}

func NewActorUsecase(ar actor.Repository, l logger.Interface) *Usecase {
	return &Usecase{ar, l}
}

func (au *Usecase) AddActor(ctx context.Context, actor *model.Actor) (uint, error) {
	id, err := au.actorRepo.AddActor(ctx, actor)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (au *Usecase) UpdateActor(ctx context.Context, actor *model.Actor) (*model.Actor, error) {
	updatedActor, err := au.actorRepo.UpdateActor(ctx, actor)
	if err != nil {
		return nil, err
	}
	return updatedActor, nil
}

func (au *Usecase) DeleteActor(ctx context.Context, actorID uint) (uint, error) {
	id, err := au.actorRepo.DeleteActor(ctx, actorID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (au *Usecase) GetActors(ctx context.Context) ([]model.ResponseActor, error) {
	actors, err := au.actorRepo.GetActors(ctx)
	if err != nil {
		return []model.ResponseActor{}, err
	}
	return actors, nil
}

func (au *Usecase) CheckActors(ctx context.Context, actors []uint) (bool, error) {
	// exist, err := au.actorRepo.CheckActors(ctx, actors)
	// if err != nil {
	// 	return false, err
	// }
	return false, nil
}

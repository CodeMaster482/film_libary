package actor

import (
	"context"
	"films_library/internal/model"
)

type (
	Usecase interface {
		AddActor(ctx context.Context, actor *model.Actor) (uint, error)
		UpdateActor(ctx context.Context, actor *model.Actor) (*model.Actor, error)
		DeleteActor(ctx context.Context, actorID uint) (uint, error)
		GetActors(ctx context.Context) ([]model.ResponseActor, error)

		CheckActors(ctx context.Context, actors []uint) (bool, error)
	}

	Repository interface {
		AddActor(ctx context.Context, actor *model.Actor) (uint, error)
		UpdateActor(ctx context.Context, actor *model.Actor) (*model.Actor, error)
		DeleteActor(ctx context.Context, actorID uint) (uint, error)
		GetActor(ctx context.Context, actorID uint) (model.Actor, error)
		GetActors(ctx context.Context) ([]model.ResponseActor, error)

		CheckActor(ctx context.Context, actor uint) (bool, error)
	}
)

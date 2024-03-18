package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	mock_actor "films_library/internal/actor/mocks"
	"films_library/internal/model"
	"films_library/pkg/logger"

	"github.com/golang/mock/gomock"
)

func TestUsecase_AddActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerMock := logger.NewMockInterface(ctrl)
	actorRepo := mock_actor.NewMockRepository(ctrl)
	usecase := NewActorUsecase(actorRepo, loggerMock)

	ctx := context.Background()

	testCases := []struct {
		name          string
		actor         *model.Actor
		expectedID    uint
		expectedError error
	}{
		{
			name:          "Valid actor",
			actor:         &model.Actor{Name: "John Doe", Sex: "M"},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name:          "Error from repository",
			actor:         &model.Actor{Name: "Invalid Actor", Sex: "F"},
			expectedID:    0,
			expectedError: errors.New("repository error"),
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actorRepo.EXPECT().AddActor(ctx, tc.actor).Return(tc.expectedID, tc.expectedError)

			id, err := usecase.AddActor(ctx, tc.actor)
			if err != tc.expectedError {
				t.Errorf("Unexpected error: expected %v, got %v", tc.expectedError, err)
			}

			if id != tc.expectedID {
				t.Errorf("Expected ID %d, got %d", tc.expectedID, id)
			}
		})
	}
}

func TestUsecase_UpdateActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerMock := logger.NewMockInterface(ctrl)
	actorRepo := mock_actor.NewMockRepository(ctrl)
	usecase := NewActorUsecase(actorRepo, loggerMock)

	ctx := context.Background()

	testCases := []struct {
		name          string
		actor         *model.Actor
		expectedActor *model.Actor
		expectedError error
	}{
		{
			name:          "Valid actor",
			actor:         &model.Actor{ID: 1, Name: "Updated John Doe", Sex: "F"},
			expectedActor: &model.Actor{ID: 1, Name: "Updated John Doe", Sex: "F"},
			expectedError: nil,
		},
		{
			name:          "Error from repository",
			actor:         &model.Actor{ID: 2, Name: "Invalid Actor", Sex: "F"},
			expectedActor: nil,
			expectedError: errors.New("repository error"),
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actorRepo.EXPECT().UpdateActor(ctx, tc.actor).Return(tc.expectedActor, tc.expectedError)

			updatedActor, err := usecase.UpdateActor(ctx, tc.actor)
			if err != tc.expectedError {
				t.Errorf("Unexpected error: expected %v, got %v", tc.expectedError, err)
			}

			if !reflect.DeepEqual(updatedActor, tc.expectedActor) {
				t.Errorf("Expected actor %v, got %v", tc.expectedActor, updatedActor)
			}
		})
	}
}

func TestUsecase_DeleteActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerMock := logger.NewMockInterface(ctrl)
	actorRepo := mock_actor.NewMockRepository(ctrl)
	usecase := NewActorUsecase(actorRepo, loggerMock)

	ctx := context.Background()

	testCases := []struct {
		name          string
		actorID       uint
		expectedID    uint
		expectedError error
	}{
		{
			name:          "Valid actor ID",
			actorID:       1,
			expectedID:    1,
			expectedError: nil,
		},
		{
			name:          "Error from repository",
			actorID:       2,
			expectedID:    0,
			expectedError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actorRepo.EXPECT().DeleteActor(ctx, tc.actorID).Return(tc.expectedID, tc.expectedError)

			id, err := usecase.DeleteActor(ctx, tc.actorID)
			if err != tc.expectedError {
				t.Errorf("Unexpected error: expected %v, got %v", tc.expectedError, err)
			}

			if id != tc.expectedID {
				t.Errorf("Expected ID %d, got %d", tc.expectedID, id)
			}
		})
	}
}

func TestUsecase_GetActors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerMock := logger.NewMockInterface(ctrl)
	actorRepo := mock_actor.NewMockRepository(ctrl)
	usecase := NewActorUsecase(actorRepo, loggerMock)

	ctx := context.Background()

	testCases := []struct {
		name          string
		actors        []model.ResponseActor
		expectedError error
	}{
		{
			name: "Valid actors",
			actors: []model.ResponseActor{
				{ActorID: 1, Name: "John Doe", Sex: "M"},
				{ActorID: 2, Name: "Jane Smith", Sex: "F"},
			},
			expectedError: nil,
		},
		{
			name:          "Error from repository",
			actors:        []model.ResponseActor{},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actorRepo.EXPECT().GetActors(ctx).Return(tc.actors, tc.expectedError)

			actors, err := usecase.GetActors(ctx)
			if err != tc.expectedError {
				t.Errorf("Unexpected error: expected %v, got %v", tc.expectedError, err)
			}

			if !reflect.DeepEqual(actors, tc.actors) {
				t.Errorf("Expected actors %v, got %v", tc.actors, actors)
			}
		})
	}
}

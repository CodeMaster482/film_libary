package model

import "time"

type Actor struct {
	ID        int    `json:"id"`
	Name      string `json:"name" validate:"required"`
	Sex       string `json:"sex" validate:"oneof=M W N"`
	BirthDate string `json:"birth_date"`
}

type ResponseActor struct {
	ActorID   uint      `json:"actor_id"`
	Name      string    `json:"name"`
	Sex       string    `json:"sex"`
	BirthDate time.Time `json:"birth_date"`
	Films     []FilmObj `json:"film"`
}

type FilmObj struct {
	Id    uint   `json:"film_id"`
	Title string `json:"title"`
}

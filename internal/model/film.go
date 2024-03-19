package model

import "time"

type Film struct {
	ID          uint64    `json:"film_id" validate:"required"`
	Title       string    `json:"title" validate:"max=150"`
	Description string    `json:"description" validate:"max=1000"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      int       `json:"rating" validate:"min=-1,max=10"`
}

type AddFilmRequest struct {
	Title       string    `json:"title" validate:"required,min=1,max=150"`
	Description string    `json:"description" validate:"max=1000"`
	ReleaseDate time.Time `json:"release_date" validate:"required"`
	Rating      int       `json:"rating" validate:"min=0,max=10"`
	Actors      []uint    `json:"actors"`
}

type FilmFilter struct {
	SortBy    string `validate:"oneof=rating release_date title"`
	SortOrder string `validate:"oneof=asc desc"`
}

// type ResponseFilm struct {
// 	MovieID     int       `json:"film_id"`
// 	Title       string    `json:"title"`
// 	Description string    `json:"description"`
// 	ReleaseDate time.Time `json:"release_date"`
// 	Rating      int       `json:"rating"`
// }

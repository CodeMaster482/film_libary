package postgresql

import (
	"context"
	"films_library/internal/model"
	"films_library/pkg/postgres"
	"fmt"
)

type Repository struct {
	db postgres.DBConn
}

func NewRepository(db postgres.DBConn) *Repository {
	return &Repository{db}
}

func (r *Repository) GetFilms(ctx context.Context, filter model.FilmFilter) ([]model.Film, error) {
	sqlQuery := `SELECT film_id, title, "description", release_date, rating FROM film`

	if filter.SortBy != "" {
		sqlQuery += " ORDER BY " + filter.SortBy
		if filter.SortOrder == "desc" {
			sqlQuery += " DESC"
		} else {
			sqlQuery += " ASC"
		}
	}

	rows, err := r.db.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []model.Film
	for rows.Next() {
		var film model.Film
		if err := rows.Scan(
			&film.ID,
			&film.Title,
			&film.Description,
			&film.ReleaseDate,
			&film.Rating,
		); err != nil {
			return nil, err
		}
		fmt.Println(film)
		films = append(films, film)
	}
	return films, nil
}

func (r *Repository) GetFilm(ctx context.Context, id uint64) (model.Film, error) {
	sqlQuery := `SELECT (film_id, title, "description", release_date, rating) FROM film WHERE film_id=$1`

	row := r.db.QueryRow(ctx, sqlQuery, id)
	var film model.Film
	err := row.Scan(&film.ID, &film.Title, &film.Description, &film.ReleaseDate, &film.Rating)
	if err != nil {
		return model.Film{}, err
	}
	return film, nil
}

func (r *Repository) AddFilm(ctx context.Context, film model.AddFilmRequest) (uint64, error) {
	sqlQuery := `INSERT INTO film (title, "description", release_date, rating) VALUES ($1, $2, $3, $4) RETURNING film_id`
	var id uint64
	err := r.db.QueryRow(ctx, sqlQuery, film.Title, film.Description, film.ReleaseDate, film.Rating).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) UpdateFilm(ctx context.Context, film model.Film) (uint64, error) {
	sqlQuery := `UPDATE film SET title=$1, "description"=$2, release_date=$3, rating=$4 WHERE film_id=$5 RETURNING film_id`

	var filmId uint64

	if err := r.db.QueryRow(ctx, sqlQuery,
		film.Title,
		film.Description,
		film.ReleaseDate,
		film.Rating,
		film.ID,
	).Scan(&filmId); err != nil {
		return 0, err
	}

	return filmId, nil
}

func (r *Repository) DeleteFilm(ctx context.Context, id uint64) (uint64, error) {
	sqlQuery := `DELETE FROM film WHERE film_id=$1`
	res, err := r.db.Exec(ctx, sqlQuery, id)
	if err != nil {
		return 0, err
	}
	rowsAffected := res.RowsAffected()
	return uint64(rowsAffected), nil

}

func (r *Repository) SearchFilm(ctx context.Context, search string) ([]model.Film, error) {
	sqlQuery := `
        SELECT DISTINCT f.film_id, f.title, f.description, f.release_date, f.rating
        FROM film f
        JOIN film_actor fa ON f.film_id = fa.film_id
        JOIN actor a ON fa.actor_id = a.actor_id
        WHERE f.title LIKE '%' || $1 || '%'
        OR a.name LIKE '%' || $1 || '%'
    `

	rows, err := r.db.Query(ctx, sqlQuery, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []model.Film
	for rows.Next() {
		var film model.Film
		if err := rows.Scan(
			&film.ID,
			&film.Title,
			&film.Description,
			&film.ReleaseDate,
			&film.Rating,
		); err != nil {
			return nil, err
		}
		films = append(films, film)
	}
	return films, nil
}

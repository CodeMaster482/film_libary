package postgres

import (
	"context"
	"films_library/internal/model"
	"films_library/pkg/postgres"
)

type Repository struct {
	db postgres.DBConn
}

func NewRepository(db postgres.DBConn) *Repository {
	return &Repository{db}
}

func (ar *Repository) AddActor(ctx context.Context, actor *model.Actor) (uint, error) {
	sqlQuery := `INSERT INTO actor (name, sex, birth_date) VALUES ($1, $2, $3) RETURNING id`

	var id uint
	if err := ar.db.QueryRow(ctx, sqlQuery,
		actor.Name,
		actor.Sex,
		actor.BirthDate,
	).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (ar *Repository) UpdateActor(ctx context.Context, actor *model.Actor) (*model.Actor, error) {
	sqlQuery := `UPDATE actor SET name = $1, sex = $2, birth_date = $3 WHERE id = $4 RETURNING id`

	var id uint
	if err := ar.db.QueryRow(ctx, sqlQuery,
		actor.Name,
		actor.Sex,
		actor.BirthDate,
		actor.ID,
	).Scan(&id); err != nil {
		return nil, &model.ErrNotFound{}
	}
	return actor, nil
}

func (ar *Repository) DeleteActor(ctx context.Context, actorID uint) (uint, error) {
	sqlQuery := `DELETE FROM actor WHERE id = $1 RETURNING id`

	var id uint
	if err := ar.db.QueryRow(ctx, sqlQuery, actorID).Scan(&id); err != nil {
		return 0, &model.ErrNotFound{}
	}
	return id, nil

}

func (ar *Repository) GetActor(ctx context.Context, actorID uint) (model.Actor, error) {
	sqlQuery := `SELECT id, name, sex, birth_date FROM actor WHERE id = $1`
	row := ar.db.QueryRow(ctx, sqlQuery, actorID)
	var actor model.Actor
	if err := row.Scan(
		&actor.ID,
		&actor.Name,
		&actor.Sex,
		&actor.BirthDate,
	); err != nil {
		return model.Actor{}, err
	}
	return actor, nil
}

func (ar *Repository) GetActors(ctx context.Context) ([]model.ResponseActor, error) {
	sqlQuery := `SELECT id, name, sex, birth_date FROM actor;`
	subQuery := ` SELECT fa.movie_id, f.title
					FROM film_actor AS fa
					JOIN film f ON f.movie_id = fa.movie_id WHERE fa.actor_id = $1;`

	rows, err := ar.db.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []model.ResponseActor
	for rows.Next() {
		var actor model.ResponseActor
		if err := rows.Scan(
			&actor.ActorID,
			&actor.Name,
			&actor.Sex,
			&actor.BirthDate,
		); err != nil {
			return nil, err
		}

		rowsFilms, err := ar.db.Query(ctx, subQuery, actor.ActorID)
		if err != nil {
			return nil, err
		}

		for rowsFilms.Next() {
			var film model.FilmObj
			if err := rowsFilms.Scan(
				&film.Id,
				&film.Title,
			); err != nil {
				return nil, err
			}
			actor.Films = append(actor.Films, film)
		}
		if err := rowsFilms.Err(); err != nil {
			return nil, err
		}

		actors = append(actors, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actors, nil
}

func (ar *Repository) CheckActor(ctx context.Context, actor uint) (bool, error) {
	return false, nil
}

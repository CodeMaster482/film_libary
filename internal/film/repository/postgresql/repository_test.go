package postgresql

import (
	"context"
	"errors"
	"films_library/internal/model"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateFilm(t *testing.T) {
	film := model.Film{
		Title:       "Updated Title",
		ReleaseDate: time.Now(),
		ID:          1,
		Description: "...",
		Rating:      8,
	}

	tests := []struct {
		name       string
		returnRows *pgxmock.Rows
		errRows    error
	}{
		{
			name:       "Success",
			returnRows: pgxmock.NewRows([]string{"film_id"}).AddRow(film.ID),
			errRows:    nil,
		},
		{
			name:       "Error",
			returnRows: &pgxmock.Rows{},
			errRows:    errors.New("mock error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mock.Close()

			repo := NewRepository(mock)

			mock.ExpectQuery(`UPDATE film SET title=$1, "description"=$2, release_date=$3, rating=$4 WHERE film_id=$5 RETURNING film_id`).
				WithArgs(film.Title, film.Description, film.ReleaseDate, film.Rating, film.ID).
				WillReturnRows(test.returnRows).
				WillReturnError(test.errRows)

			updatedFilmID, err := repo.UpdateFilm(context.Background(), film)

			if test.errRows != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, film.ID, updatedFilmID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

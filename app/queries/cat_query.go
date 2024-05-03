package queries

import (
	"github.com/iwinardhyas/sprint_backend/app/models"
	"github.com/jmoiron/sqlx"
)

type CatQueries struct {
	*sqlx.DB
}

func (q *CatQueries) CreateCat(u *models.Cat) error {

	query := `INSERT INTO cats VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := q.Exec(
		query,
		u.ID, u.Name, u.Race, u.Sex, u.AgeInMonth, u.Description, u.ImageUrls,
	)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil

}

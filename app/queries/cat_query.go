package queries

import (
	"fmt"
	"strings"
	"time"

	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/jmoiron/sqlx"
)

type CatQueries struct {
	*sqlx.DB
}

func (q *CatQueries) GetCats(appendQuery string) ([]models.CatData, error) {
	query := "SELECT id,name,race,sex,age_in_month,image_urls,description,has_matched,created_at FROM cats WHERE cat_status=1"
	rows, err := q.Query(query + appendQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.CatData
	for rows.Next() {
		var each = models.CatData{}
		var imgUrlStr string
		var err = rows.Scan(
			&each.ID,
			&each.Name,
			&each.Race,
			&each.Sex,
			&each.AgeInMonth,
			&imgUrlStr,
			&each.Description,
			&each.HasMatched,
			&each.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		each.ImageUrls = strings.Split(imgUrlStr, ",")
		result = append(result, each)
	}
	return result, nil
}

func (q *CatQueries) CreateCat(c *models.Cat) error {
	query := `INSERT INTO cats (id, user_id, name, race, sex, age_in_month, description, image_urls, has_matched, created_at, updated_at)
           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := q.Exec(
		query,
		c.ID,
		c.UserID,
		c.Name,
		c.Race,
		c.Sex,
		c.AgeInMonth,
		c.Description,
		c.ImageUrls,
		c.HasMatched,
		c.CreatedAt,
		c.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (q *CatQueries) DeleteCat(catId string, userID string) error {
	query := `UPDATE cats SET deleted_at = $1 WHERE id = $2 AND user_id = $3`

	deletedAt := time.Now()
	fmt.Printf(catId, userID)
	res, err := q.Exec(query, deletedAt, catId, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no id found")
	}

	return nil
}

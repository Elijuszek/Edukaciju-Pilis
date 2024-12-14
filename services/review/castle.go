package review

// TODO: Use sqlx
import (
	"database/sql"
	"educations-castle/types"
	"fmt"
)

type Castle struct {
	db *sql.DB
}

func NewCastle(db *sql.DB) *Castle {
	return &Castle{db: db}
}

func scanRowIntoReview(rows *sql.Rows) (*types.Review, error) {
	r := new(types.Review)

	err := rows.Scan(
		&r.ID,
		&r.Date,
		&r.Comment,
		&r.Rating,
		&r.FkUserID,
		&r.FkActivityID,
	)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Castle) ListReviews() ([]*types.Review, error) {
	rows, err := c.db.Query("SELECT * FROM review")
	if err != nil {
		return nil, err
	}

	var reviews []*types.Review

	for rows.Next() {
		r := new(types.Review)
		r, err = scanRowIntoReview(rows)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (c *Castle) GetReviewByID(id int) (*types.Review, error) {
	rows, err := c.db.Query("SELECT * FROM review WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	r := new(types.Review)
	for rows.Next() {
		r, err = scanRowIntoReview(rows)
		if err != nil {
			return nil, err
		}
	}

	if r.ID == 0 {
		return nil, sql.ErrNoRows
	}

	return r, nil
}

func (c *Castle) GetReviewFromActivityByID(idActivity int, idUser int) (*types.Review, error) {
	rows, err := c.db.Query("SELECT * FROM review WHERE fk_Activityid = ? AND fk_Userid = ?", idActivity, idUser)
	if err != nil {
		return nil, err
	}
	r := new(types.Review)
	for rows.Next() {
		r, err = scanRowIntoReview(rows)
		if err != nil {
			return nil, err
		}
	}

	if r.ID == 0 {
		return nil, fmt.Errorf("review not found")
	}

	return r, nil
}

func (c *Castle) CreateReview(review types.Review) error {
	_, err := c.db.Exec("INSERT INTO review (comment, rating, fk_Userid, fk_Activityid) VALUES (?,?,?,?,?)",
		review.Comment, review.Rating, review.FkUserID, review.FkActivityID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Castle) UpdateReview(review types.Review) error {
	_, err := c.db.Exec(
		"UPDATE review SET comment = ?, rating = ?, fk_Userid = ?, fk_Activityid = ? WHERE id = ?",
		review.Comment, review.Rating, review.FkUserID, review.FkActivityID, review.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Castle) DeleteReviewByID(id int) error {
	_, err := c.db.Exec(
		"DELETE FROM review WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Castle) ListReviewsFromPackage(id int) ([]*types.Review, error) {
	rows, err := c.db.Query(`
		SELECT review.*
		FROM review
		JOIN activity ON review.fk_Activityid = activity.id
		JOIN package ON activity.fk_Packageid = package.id
		WHERE package.id = ?
	`, id)

	var reviews []*types.Review

	for rows.Next() {
		r := new(types.Review)
		r, err = scanRowIntoReview(rows)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

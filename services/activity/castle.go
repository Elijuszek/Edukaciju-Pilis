package activity

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

func scanRowIntoActivity(rows *sql.Rows) (*types.Activity, error) {
	activity := new(types.Activity)

	err := rows.Scan(
		&activity.ID,
		&activity.Name,
		&activity.Description,
		&activity.BasePrice,
		&activity.CreationDate,
		&activity.Hidden,
		&activity.Verified,
		&activity.Category,
		&activity.AverageRating,
		&activity.FkPackageID,
	)

	if err != nil {
		return nil, err
	}

	return activity, nil
}

func scanRowIntoPackage(rows *sql.Rows) (*types.Package, error) {
	p := new(types.Package)

	err := rows.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.FkOrganizerID,
	)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (c *Castle) ListActivities() ([]*types.Activity, error) {
	rows, err := c.db.Query("SELECT * FROM activity")
	if err != nil {
		return nil, err
	}

	var activity []*types.Activity

	for rows.Next() {
		r := new(types.Activity)
		r, err = scanRowIntoActivity(rows)
		if err != nil {
			return nil, err
		}
		activity = append(activity, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return activity, nil
}

func (c *Castle) CreateActivity(activity types.Activity) error {
	var categoryID int
	err := c.db.QueryRow("SELECT id_Category FROM category WHERE name = ?", activity.Category).Scan(&categoryID)
	if err != nil {
		return fmt.Errorf("failed to find category '%s': %v", activity.Category, err)
	}

	_, err = c.db.Exec(
		"INSERT INTO activity (name, description, basePrice, hidden, category, fk_Packageid) VALUES (?,?,?,?,?,?)",
		activity.Name, activity.Description, activity.BasePrice, activity.Hidden, categoryID, activity.FkPackageID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Castle) GetActivityByID(id int) (*types.Activity, error) {
	rows, err := c.db.Query("SELECT * FROM activity WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	a := new(types.Activity)
	for rows.Next() {
		a, err = scanRowIntoActivity(rows)
		if err != nil {
			return nil, err
		}
	}

	if a.ID == 0 {
		return nil, sql.ErrNoRows
	}

	return a, nil
}

func (c *Castle) GetActivityInsidePackageByName(activityName string, packageID int) (*types.Activity, error) {
	rows, err := c.db.Query("SELECT * FROM activity WHERE name = ? AND fk_Packageid = ?", activityName, packageID)
	if err != nil {
		return nil, err
	}

	a := new(types.Activity)
	for rows.Next() {
		a, err = scanRowIntoActivity(rows)
		if err != nil {
			return nil, err
		}
	}
	if a.ID == 0 {
		return nil, fmt.Errorf("activity not found")
	}

	return a, nil
}

func (c *Castle) UpdateActivity(activity types.Activity) error {
	var categoryID int
	err := c.db.QueryRow("SELECT id_Category FROM category WHERE name = ?", activity.Category).Scan(&categoryID)
	if err != nil {
		return fmt.Errorf("failed to find category '%s': %v", activity.Category, err)
	}
	_, err = c.db.Exec(
		`UPDATE activity 
		SET name = ?, description = ?, basePrice = ?, hidden = ?, category = ?, fk_Packageid = ? 
		WHERE id = ?`,
		activity.Name, activity.Description, activity.BasePrice, activity.Hidden, categoryID,
		activity.FkPackageID, activity.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Castle) DeleteActivity(id int) error {
	_, err := c.db.Exec(
		"DELETE FROM activity WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Castle) CreatePackage(p types.Package) error {
	_, err := c.db.Exec(
		"INSERT INTO package (name, description, price, fk_Organizerid) VALUES (?,?,?,?)",
		p.Name, p.Description, p.Price, p.FkOrganizerID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Castle) GetPackageByName(name string) (*types.Package, error) {
	rows, err := c.db.Query("SELECT * FROM package WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	p := new(types.Package)
	for rows.Next() {
		p, err = scanRowIntoPackage(rows)
		if err != nil {
			return nil, err
		}
	}

	if p.ID == 0 {
		return nil, fmt.Errorf("package not found")
	}

	return p, nil
}

func (c *Castle) DeletePackage(id int) error {
	_, err := c.db.Exec(
		"DELETE FROM package WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

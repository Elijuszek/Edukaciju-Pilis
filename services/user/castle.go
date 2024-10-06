package user

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

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.RegistrationDate,
		&user.LastLoginDate,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Castle) GetUserByID(id int) (*types.User, error) {
	rows, err := c.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Castle) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM user WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil

}

func (c *Castle) GetUserByUsername(username string) (*types.User, error) {
	rows, err := c.db.Query("SELECT * FROM user WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (c *Castle) CreateUser(user types.User) error {
	_, err := c.db.Exec("INSERT INTO user (username, password, email) VALUES (?,?,?)", user.Username,
		user.Password, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (c *Castle) DeleteUser(id int) error {
	_, err := c.db.Exec(
		"DELETE FROM user WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Castle) CreateOrganizer(organizer types.Organizer) error {
	_, err := c.db.Exec("INSERT INTO organizer (id, description) VALUES (?,?)", organizer.ID,
		organizer.Description)
	if err != nil {
		return err
	}

	return nil
}

func (c *Castle) DeleteOrganizer(id int) error {
	_, err := c.db.Exec(
		"DELETE FROM organizer WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Castle) DeleteAdministrator(id int) error {
	_, err := c.db.Exec(
		"DELETE FROM administrator WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

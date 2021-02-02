package mysql

import (
	"database/sql"
	"errors"
	"github.com/aovlllo/vue-template/pkg/model"
)

var ErrAlreadyExist = errors.New("user already exists")

// CreateUser creates a new user
func (db *DB) CreateUser(u *model.User) error {
	oldU, err := db.GetUserByEmail(u.Email)
	if err != nil {
		return err
	} else if oldU != nil {
		return ErrAlreadyExist
	}

	res, err := db.client.Exec(
		"INSERT INTO users (name, email, password, second_name, birth, sex, city, interests) VALUES(?,?,?,?,?,?,?,?)",
		u.Name, u.Email, u.Password, u.SecondName, u.Birth, u.Sex, u.City, u.Interests,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	u.ID = int(id)
	return nil
}

// GetUser returns a user
func (db *DB) GetUser(id int) (*model.User, error) {
	row := db.client.QueryRow("SELECT id, name, email, password, second_name, birth, sex, city, interests FROM users WHERE id = ?", id)
	u := model.User{}

	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.SecondName, &u.Birth, &u.Sex, &u.City, &u.Interests); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

// GetUserByEmail returns a user by his email address
func (db *DB) GetUserByEmail(email string) (*model.User, error) {
	row := db.client.QueryRow("SELECT id, name, email, password, second_name, birth, sex, city, interests FROM users WHERE email = ?", email)
	u := model.User{}

	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.SecondName, &u.Birth, &u.Sex, &u.City, &u.Interests); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

// UpdateUser saves the given user struct
func (db *DB) UpdateUser(u *model.User) error {
	if _, err := db.client.Exec(
		"UPDATE users SET name = ?, email = ?, second_name = ?, birth = ?, sex = ?, city = ?, interests = ? WHERE id = ?",
		u.Name, u.Email, u.SecondName, u.Birth, u.Sex, u.City, u.Interests, u.ID,
	); err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes the user with the given id
func (db *DB) DeleteUser(id int) error {
	if _, err := db.client.Exec("DELETE FROM users WHERE id = ?", id); err != nil {
		return err
	}
	return nil
}

package database

import (
	"CourseProject/auth_service/internal/entity"
	"database/sql"
	"errors"
	"fmt"
)

type UserStorage struct {
	DB *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{DB: db}
}

func (us *UserStorage) Post(newUser entity.User) error {
	query := "INSERT INTO users (id, username, password) VALUES ($1, $2, $3)"
	result, err := us.DB.Exec(query, newUser.ID, newUser.Username, newUser.Password)
	if err != nil {
		return fmt.Errorf("unable to insert new user into users database: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("zero rows were inserted")
	}
	return nil
}

func (us *UserStorage) Get(userName string) (*entity.User, error) {
	query := "SELECT id, username, password FROM users WHERE username = $1"

	row := us.DB.QueryRow(query, userName)
	var userFromDB entity.User
	err := row.Scan(&userFromDB.ID, &userFromDB.Username, &userFromDB.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("there is no rows with passed username: %w", err)
		}
		return nil, fmt.Errorf("unable to scan user from selected row: %w", err)
	}
	return &userFromDB, nil
}

package database

import (
	"CourseProject/auth_service/internal/entity"
	"database/sql"
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

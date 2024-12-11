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
	query := "INSERT INTO users (id, username, email, phone, password) VALUES ($1, $2, $3, CASE WHEN $4 = '' THEN NULL ELSE $4 END, $5)"
	result, err := us.DB.Exec(query, newUser.ID, newUser.Username, newUser.Email, newUser.Phone, newUser.Password)
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

func (us *UserStorage) GetUser(query string, args ...interface{}) (*entity.User, error) {
	row := us.DB.QueryRow(query, args...)
	var userFromDB entity.User
	err := row.Scan(&userFromDB.ID, &userFromDB.Username, &userFromDB.Email, &userFromDB.Phone, &userFromDB.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("there is no rows with passed username: %w", err)
		}
		return nil, fmt.Errorf("unable to scan user from selected row: %w", err)
	}
	return &userFromDB, nil
}

func (us *UserStorage) GetByEmail(email string) (*entity.User, error) {
	query := "SELECT id, username, email, phone, password FROM users WHERE email = $1"
	return us.GetUser(query, email)
}

func (us *UserStorage) GetByName(userName string) (*entity.User, error) {
	query := "SELECT id, username, email, phone, password FROM users WHERE username = $1"
	return us.GetUser(query, userName)
}

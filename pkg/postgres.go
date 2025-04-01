package pkg

import (
	"CourseProject/auth_service/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func PostgresConnect(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %s", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping the database: %s", err)
	}
	return db, nil
}

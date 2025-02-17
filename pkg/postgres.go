package pkg

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func PostgresConnect(dbName string) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("host=postgres_users port=5432 user=admin password=admin dbname=%s sslmode=disable", dbName)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %s", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping the database: %s", err)
	}
	return db, nil
}

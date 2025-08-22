package dbkit

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Open(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

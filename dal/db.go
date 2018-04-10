package dal

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//StartDB opens a connection to the given sql database.
func StartDB(dbAddress string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbAddress)
	if err != nil {
		return nil, err
	}

	return db, nil
}

package dal

import (
	"context"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
)

//StartDB opens a connection to psql database.
func StartDB() (conn *pgx.Conn) {
	connString := os.Getenv("PSQL_URL")

	m, err := migrate.New(
		"file://dal/migrations",
		connString)

	if err != nil {
		log.Fatalf("Unable to connect w/ migrate to db. Shutting down. Err: %v\n", err)
	}

	migrateErr := m.Up()

	if migrateErr != nil && migrateErr != migrate.ErrNoChange {
		log.Fatalf("Unable to migrate to db. Shutting down. Err: %v\n", migrateErr)
	}

	conn, dbErr := pgx.Connect(context.Background(), connString)

	if dbErr != nil {
		log.Fatalf("Unable to connect to DB. Shutting down. Err: %v\n", dbErr)
	}
	return conn
}

//StartDB opens a connection to the given sql database.
// func StartDB(dbAddress string) (*sql.DB, error) {
// 	db, err := sql.Open("mysql", dbAddress)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }

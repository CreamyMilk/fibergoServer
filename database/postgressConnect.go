package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var PDB *sql.DB

// Database settings
const (
	host     = "localhost"
	port     = 5432 // Default port
	user     = "postgres"
	password = "password"
	dbname   = "fiber_demo"
)

func PGConnect() error {
	var err error
	PDB, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		return err
	}
	if err = PDB.Ping(); err != nil {
		return err
	}
	return nil
}

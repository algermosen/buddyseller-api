package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func InitDB() {
	url := getConnectionString()
	fmt.Println(url)
	var err error
	DB, err = sql.Open("pgx", url)

	if err != nil {
		panic("Could not connect to database")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	createTables()
}

func createTables() {
	var err error
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id serial PRIMARY KEY ,
		name character(30) NOT NULL,
		description text NOT NULL,
		location text NOT NULL,
		datetime timestamp NOT NULL,
		user_id integer
	)
	`

	_, err = DB.Exec(createEventsTable)

	if err != nil {
		panic(err.Error())
	}

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY ,
		email character(30) NOT NULL UNIQUE,
		password text NOT NULL
	)
	`

	_, err = DB.Exec(createUsersTable)

	if err != nil {
		panic("Could not create table: Users")
	}

	createUsersEventsRelation := `
	DO $$
	BEGIN

		BEGIN
			ALTER TABLE events
			ADD CONSTRAINT fk_events_users FOREIGN KEY (user_id) REFERENCES users(id);
		EXCEPTION WHEN duplicate_object THEN
		END;

	END $$;
	`

	_, err = DB.Exec(createUsersEventsRelation)

	if err != nil {
		panic("Could not reference users table in events")
	}

	createRegistrationsTable := `
	CREATE TABLE IF NOT EXISTS registrations (
		id serial PRIMARY KEY,
		user_id integer,
		event_id integer,
		FOREIGN KEY (event_id) REFERENCES events(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	)
	`

	_, err = DB.Exec(createRegistrationsTable)

	if err != nil {
		panic("Could not create registrations table")
	}
}

func getConnectionString() string {
	var user, pass, host, port, dbName string

	user = os.Getenv("PG_DB_USER")
	pass = os.Getenv("PG_DB_PASS")
	host = os.Getenv("PG_DB_HOST")
	port = os.Getenv("PG_DB_PORT")
	dbName = os.Getenv("PG_DB_NAME")

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbName)
}

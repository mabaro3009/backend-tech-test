package pg

import (
	"database/sql"
	"fmt"
	"log"

	"reby/app/config"
)

func InitDB(conf *config.Config) *sql.DB {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.DBHost,
		conf.DBPort,
		conf.DBUser,
		conf.DBPassword,
		conf.DBName,
	)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	initTables(db)
	return db
}

func initTables(db *sql.DB) {
	userTable := `CREATE TABLE IF NOT EXISTS "user" (id varchar(255) PRIMARY KEY);`
	if _, err := db.Exec(userTable); err != nil {
		log.Fatal(err)
	}

	vehicleTable := `CREATE TABLE IF NOT EXISTS "vehicle" (id varchar(255) PRIMARY KEY);`
	if _, err := db.Exec(vehicleTable); err != nil {
		log.Fatal(err)
	}

	rideTable :=
		`CREATE TABLE IF NOT EXISTS "ride" (
	id varchar(255) PRIMARY KEY,
	vehicle_id varchar(255) REFERENCES vehicle(id),
	user_id varchar(255) REFERENCES "user"(id),
	started_at TIMESTAMP NOT NULL,
	finished_at TIMESTAMP,
    price_value int,
    price_currency varchar(255)
);`
	if _, err := db.Exec(rideTable); err != nil {
		log.Fatal(err)
	}
}

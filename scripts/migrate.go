package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var (
	host           string = os.Getenv("DB_HOST")
	port           string = os.Getenv("DB_PORT")
	databaseName   string = os.Getenv("DB_DATABASE")
	dbUsername     string = os.Getenv("DB_USERNAME")
	dbPassword     string = os.Getenv("DB_PASSWORD")
	httpListenAddr string = ":3000"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s sslmode=disable",
		host, port, dbUsername, dbPassword)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatal("db: ", err)
	}
	defer db.Close()

	log.Println(createDatabaseIfNotExists(db, databaseName))

	psqlInfo = fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, dbUsername, dbPassword, databaseName)
	db, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatal("db: ", err)
	}
	defer db.Close()

	log.Fatal(createTableUsersIfNotExists(db))
}

func createDatabaseIfNotExists(db *sqlx.DB, dbName string) error {
	// Check if the database already exists
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName)
	if err != nil {
		return err
	}

	if !exists {
		// Create the database if it does not exist
		_, err := db.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			return err
		}
	}

	return nil
}

func createTableUsersIfNotExists(db *sqlx.DB) error {
	qu := `CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	ALTER TABLE users
    ADD CONSTRAINT unique_email UNIQUE (email);
	`
	_, err := db.Exec(qu)
	return err
}

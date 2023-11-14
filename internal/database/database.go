package database

import (
	"fmt"
	"log"
	"os"

	"github.com/abdoroot/authentication-service/internal/auth"
	pb "github.com/abdoroot/authentication-service/proto"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var host string = os.Getenv("DB_HOST")
var port string = os.Getenv("DB_PORT")
var databaseName string = os.Getenv("DB_DATABASE")
var dbUsername string = os.Getenv("DB_USERNAME")
var dbPassword string = os.Getenv("DB_PASSWORD")

type DB struct {
	db *sqlx.DB
}

func NewDB() (*DB, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	return &DB{
		db: db,
	}, nil
}

func Connect() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, dbUsername, dbPassword, databaseName)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (in *DB) Insert(req *pb.SignUpRequest) error {
	var err error
	req.Password, err = auth.GetPasswordHash(req.Password)
	if err != nil {
		return err
	}
	_, err = in.db.NamedExec(`INSERT INTO users (name, email, password,created_at)
        VALUES (:name, :email,:password,:createdAt)`, req)
	if err != nil {
		log.Panicln(req)
		return err
	}
	return nil
}

func (in *DB) Migrate() error {
	qu := `CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := in.db.Exec(qu)
	return err
}

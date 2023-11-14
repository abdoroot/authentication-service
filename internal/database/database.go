package database

import (
	"fmt"
	"log"
	"os"

	pb "github.com/abdoroot/authentication-service/proto"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var host string = os.Getenv("DB_HOST")
var port string = os.Getenv("DB_PORT")
var databaseName string = os.Getenv("DB_DATABASE")
var dbUsername string = os.Getenv("DB_USERNAME")
var dbPassword string = os.Getenv("DB_PASSWORD")

// map use to
var dataMp map[string]any

type DB struct {
	db *sqlx.DB
}

type loginRequest struct {
	Email    string `db:"email"`
	Password string `db:"password"`
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

	dataMp = make(map[string]any)

	dataMp["name"] = req.Name
	dataMp["email"] = req.Email
	dataMp["password"], err = GetPasswordHash(req.Password)
	if err != nil {
		log.Println("Err Hashing :", err)
		return err
	}

	//insert data to db
	_, err = in.db.NamedExec(`INSERT INTO users (name, email, password)
        VALUES (:name,:email,:password)`, dataMp)
	if err != nil {
		log.Println("Err Insert data to db :", req, err)
		return err
	}
	return nil
}

func (in *DB) Login(req *pb.LoginRequest) (string, error) {
	lg := loginRequest{}
	err := in.db.Get(&lg, `select email,password from users where email=$1`, req.Email)

	if err != nil {
		log.Println(lg)
		return "", err
	}
	//Compare hash with the password
	if IsHashEqPass(lg.Password, req.Password) {
		return "genratedJwtToke", nil
	}
	return "", fmt.Errorf("some thing went wrong")
}

func (in *DB) Migrate() error {
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
	_, err := in.db.Exec(qu)
	return err
}

func GetPasswordHash(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hash), err
}

func IsHashEqPass(hash string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

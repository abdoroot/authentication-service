package auth

import (
	"fmt"
	"log"
	"os"

	pb "github.com/abdoroot/authentication-service/proto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var host string = os.Getenv("DB_HOST")
var port string = os.Getenv("DB_PORT")
var databaseName string = os.Getenv("DB_DATABASE")
var dbUsername string = os.Getenv("DB_USERNAME")
var dbPassword string = os.Getenv("DB_PASSWORD")

// glopal DB
var gdb *sqlx.DB

// map use to
var dataMp map[string]any

type DB struct {
	db *sqlx.DB
}

type User struct {
	UserId   string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type GetProfileResponse struct {
	Email string `db:"email"`
	Name  string `db:"name"`
}

func NewDB() (*DB, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}

	//set the global db
	gdb = db

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

func (in *DB) Login(req *pb.LoginRequest) (map[string]string, error) {
	usr := User{}
	err := in.db.Get(&usr, `select id,email,password from users where email=$1`, req.Email)

	if err != nil {
		log.Println(usr)
		return nil, err
	}
	//Compare hash with the password
	if IsHashEqPass(usr.Password, req.Password) {
		return GenerateToken(usr.UserId, usr.Email)
	}
	return nil, fmt.Errorf("some thing went wrong")
}

func (in *DB) Update(req *pb.UpdateRequest, claims jwt.MapClaims) error {
	userId := claims["user_id"].(string)
	name := req.Name
	if req.Password != "" {
		password, err := GetPasswordHash(req.Password)
		if err != nil {
			return err
		}
		_, err = in.db.Exec(`update users set name=$1,password=$2 where id=$3`, name, password, userId)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		//password empty
		log.Println("password empty")
		_, err := in.db.Exec(`update users set name=$1 where id=$2`, name, userId)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (in *DB) GetProfile(claims jwt.MapClaims) (*GetProfileResponse, error) {
	userId := claims["user_id"].(string)
	gr := &GetProfileResponse{}
	err := in.db.Get(gr, `select name,email from users where id=$1`, userId)
	if err != nil {
		return nil, err
	}
	return gr, nil
}

func FindUserById(id string) (*User, error) {
	usr := &User{}
	err := gdb.Get(usr, `select id,email from users where id=$1`, id)

	if err != nil {
		log.Println(usr)
		return nil, err
	}

	return usr, nil
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

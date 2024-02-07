package store

import (
	"context"
	"fmt"
	"log"

	"github.com/abdoroot/authentication-service/internal/types"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Storer interface {
	CreateUser(context.Context, *types.User) (*types.User, error)
	UpdateUser(context.Context, *types.User) (*types.User, error)
	GetUserByEmailPassword(context.Context, string, string) (*types.User, error)
	GetUserById(context.Context, string) (*types.User, error)
}

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	resp, err := s.db.NamedExec(`INSERT INTO users (name, email, password)
        VALUES (:name,:email,:password)`, user)
	if err != nil {
		log.Println("err Insert data to db:", user, err)
		return user, err
	}
	insetedUserId, _ := resp.LastInsertId()
	user.ID = int(insetedUserId)
	return user, nil
}

func (s *UserStore) UpdateUser(ctx context.Context, user *types.User) (*types.User, error) {
	var err error
	if user.Password != "" {
		_, err = s.db.Exec(`update users set name=$1,password=$2 where id=$3`, user.Name, user.Password, user.ID)
		if err != nil {
			log.Println("update user err", err)
			return user, err
		}
	} else {
		//password empty
		_, err := s.db.Exec(`update users set name=$1 where id=$2`, user.Name, user.ID)
		if err != nil {
			log.Println("update user err", err)
			return user, err
		}
	}
	return user, nil
}

func (s *UserStore) GetUserByEmailPassword(ctx context.Context, email, password string) (*types.User, error) {
	user := &types.User{}
	err := s.db.Get(user, `select * from users where email=$1`, email)
	if err != nil {
		return nil, err
	}

	//Compare hash with the password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("error email or password")
	}

	return user, nil
}

func (s *UserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	user := &types.User{}
	err := s.db.Get(user, `select * from users where id=$1`, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

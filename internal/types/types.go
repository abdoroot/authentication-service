package types

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type HttpApiFunc func(http.ResponseWriter, *http.Request) error

type User struct {
	ID        int    `json:"id,omitempty" db:"id"`
	Name      string `json:"name" db:"name"`
	Email     string `json:"email,omitempty" db:"email"`
	Password  string `json:"-" db:"password"`
	CreatedAt string `json:"created_at,omitempty" db:"created_at"`
}

type LoginParam struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserParam struct {
	Name     string
	Email    string
	Password string
}

func (param CreateUserParam) CreateUserFromParam() (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		Name:     param.Name,
		Email:    param.Email,
		Password: string(hash),
	}, nil
}

type GetUserParam struct {
	Email    string
	Password string
}

type UpdateUserParam struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (param UpdateUserParam) CreateUpdateRequest() (*User, error) {
	user := &User{}
	if param.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hash)
	}

	user.Name = param.Name
	return user, nil
}

package auth

import (
	"context"

	"github.com/abdoroot/authentication-service/internal/store"
	"github.com/abdoroot/authentication-service/internal/types"
)

type Auth struct {
	Store store.Storer
}

func NewAuth(store store.Storer) *Auth {
	return &Auth{
		Store: store,
	}
}

func (a *Auth) SignUp(ctx context.Context, user *types.User) (*types.User, error) {
	createUser, err := a.Store.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return createUser, nil
}

func (a *Auth) Login(ctx context.Context, param *types.LoginParam) (*types.User, error) {
	user, err := a.Store.GetUserByEmailPassword(ctx, param.Email, param.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *Auth) Update(ctx context.Context, user *types.User) (*types.User, error) {
	//update opration
	//todo validate inputs
	user, err := a.Store.UpdateUser(context.Background(), user)
	if err != nil {
		return nil, err
	}
	//updated
	return user, nil
}

func (a Auth) UserProfile(ctx context.Context, id string) (*types.User, error) {
	//update opration
	//todo validate inputs
	user, err := a.Store.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	//retrive
	return user, nil
}

package user

import (
	"api/pkg/model"
	"context"
)

type RepoUser interface {
	GetUser(ctx context.Context, userData model.User) (*model.User, error)
	CreateUser(ctx context.Context, userData model.User) (*model.User, error)
}

type ControllerUser struct {
	repoUser RepoUser
}

func New(repoUser RepoUser) *ControllerUser {
	return &ControllerUser{repoUser: repoUser}
}

func (u *ControllerUser) GetUser(ctx context.Context, userData model.User) (*model.User, error) {
	response, err := u.repoUser.GetUser(ctx, userData)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (u *ControllerUser) CreateUser(ctx context.Context, userData model.User) (*model.User, error) {

	response, err := u.repoUser.CreateUser(ctx, userData)

	if err != nil {
		return nil, err
	}
	return response, nil
}

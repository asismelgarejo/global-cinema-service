package user

import (
	"api/pkg/model"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrUserNotFound = errors.New("invalid username or password")

type RepoUser struct {
	collection *mongo.Collection
}

func New(collection *mongo.Collection) *RepoUser {
	return &RepoUser{collection: collection}
}

func (u *RepoUser) CreateUser(ctx context.Context, userData model.User) (*model.User, error) {

	user, err := u.GetUser(ctx, userData)
	if errors.Is(err, ErrUserNotFound) {
	} else if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, fmt.Errorf("username %v already exists", user.Username)
	}

	userData.Id = primitive.NewObjectID()
	h := sha256.New()
	userData.Password = fmt.Sprintf("%x", h.Sum([]byte(userData.Password)))

	_, err = u.collection.InsertOne(ctx, userData)
	if err != nil {
		return nil, err
	}
	userData.Password = ""
	return &userData, nil

}

func (u *RepoUser) GetUser(ctx context.Context, userData model.User) (*model.User, error) {
	var user model.User
	h := sha256.New()

	cur := u.collection.FindOne(ctx, bson.M{
		"username": userData.Username,
		"password": fmt.Sprintf("%x", h.Sum([]byte(userData.Password))),
	})
	if cur.Err() != nil {
		return nil, ErrUserNotFound
	}

	if err := cur.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

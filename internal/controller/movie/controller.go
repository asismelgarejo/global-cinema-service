package controller

import (
	"api/pkg/model"
	"context"
)

type MovieRepository interface {
	CreateMovie(ctx context.Context, movieData model.Movie) (*model.Movie, error)
	GetMovies(ctx context.Context) ([]*model.Movie, error)
	UpdateMovie(ctx context.Context, id string, movieData model.Movie) (*model.Movie, error)
	DeleteMovie(ctx context.Context, id string) (string, error)
	GetMovie(ctx context.Context, id string) (*model.Movie, error)
}

type MovieController struct {
	repo MovieRepository
}

func New(repo MovieRepository) *MovieController {
	return &MovieController{repo: repo}
}

func (c *MovieController) CreateMovie(ctx context.Context, movieData model.Movie) (*model.Movie, error) {
	response, err := c.repo.CreateMovie(ctx, movieData)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (c *MovieController) GetMovies(ctx context.Context) ([]*model.Movie, error) {
	movieData, err := c.repo.GetMovies(ctx)
	if err != nil {
		return nil, err
	}
	return movieData, nil
}
func (c *MovieController) UpdateMovie(ctx context.Context, id string, movieData model.Movie) (*model.Movie, error) {
	response, err := c.repo.UpdateMovie(ctx, id, movieData)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (c *MovieController) DeleteMovie(ctx context.Context, id string) (string, error) {
	response, err := c.repo.DeleteMovie(ctx, id)
	if err != nil {
		return "", err
	}
	return response, nil
}
func (c *MovieController) GetMovie(ctx context.Context, id string) (*model.Movie, error) {
	response, err := c.repo.GetMovie(ctx, id)
	if err != nil {
		return nil, err
	}
	return response, nil
}

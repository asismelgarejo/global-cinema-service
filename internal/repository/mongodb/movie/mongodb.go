package movie

import (
	"api/pkg/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepoMovie struct {
	collection *mongo.Collection
}

func New(collection *mongo.Collection) *RepoMovie {
	return &RepoMovie{collection: collection}
}

func (m *RepoMovie) CreateMovie(ctx context.Context, movieData model.Movie) (*model.Movie, error) {
	movieData.Id = primitive.NewObjectID()
	_, err := m.collection.InsertOne(ctx, movieData)
	if err != nil {
		return nil, err
	}
	return &movieData, nil
}
func (m *RepoMovie) GetMovies(ctx context.Context) ([]*model.Movie, error) {
	cur, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	movies := []*model.Movie{}
	for cur.Next(ctx) {
		var movie model.Movie
		cur.Decode(&movie)
		movies = append(movies, &movie)
	}

	return movies, nil
}
func (m *RepoMovie) UpdateMovie(ctx context.Context, id string, movieData model.Movie) (*model.Movie, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	_, err = m.collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.D{
		{"$set", bson.D{
			{"title", movieData.Title},
			{"genres", movieData.Genres},
			{"year", movieData.Year},
			{"directors", movieData.Directors},
			{"synopsis", movieData.Synopsis},
			{"poster", movieData.Poster},
			{"country", movieData.Country},
		}},
	})
	if err != nil {
		return nil, err
	}
	movieData.Id = objectId
	return &movieData, nil
}
func (m *RepoMovie) DeleteMovie(ctx context.Context, id string) (string, error) {

	_, err := m.GetMovie(ctx, id)
	if err != nil {
		return "", err
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	filter := bson.M{"_id": objectId}

	_, err = m.collection.DeleteOne(ctx, filter)
	if err != nil {
		return "", err
	}
	return id, nil
}
func (m *RepoMovie) GetMovie(ctx context.Context, id string) (*model.Movie, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	response := m.collection.FindOne(ctx, bson.M{"_id": objectId})
	if response.Err() != nil {
		if response.Err() == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("movie not found")
		}
		return nil, response.Err()
	}

	var movie model.Movie
	if err := response.Decode(&movie); err != nil {
		return nil, err
	}
	return &movie, nil
}

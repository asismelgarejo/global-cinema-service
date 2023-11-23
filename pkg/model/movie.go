package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Movie struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Title     string             `json:"title" bson:"title"`
	Genres    []string           `json:"genres" bson:"genres"`
	Year      int                `json:"year" bson:"year"`
	Directors []string           `json:"directors" bson:"directors"`
	Synopsis  string             `json:"synopsis" bson:"synopsis"`
	Poster    string             `json:"poster" bson:"poster"`
	Country   string             `json:"country" bson:"country"`
}

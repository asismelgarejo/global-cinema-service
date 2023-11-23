package main

import (
	ctrlMovie "api/internal/controller/movie"
	handlerMovie "api/internal/handler/http/movie"
	repoMovie "api/internal/repository/mongodb/movie"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/yaml.v2"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var database *mongo.Database
var redisClient *redis.Client

func init() {
	//---Reading config file
	configFile := "./configs/base_dev.yaml"
	if os.Getenv("mode") == "production" {
		configFile = "./configs/base_prod.yaml"
	}
	f, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var ServiceConfig serviceConfig

	if err := yaml.NewDecoder(f).Decode(&ServiceConfig); err != nil {
		panic(err)
	}
	//---Reading config file

	//---Database
	ctx := context.Background()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(ServiceConfig.DBConfig.StrConn))
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	database = client.Database(ServiceConfig.DBConfig.DBName)
	//---Database

	//---Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     ServiceConfig.RedisConfig.StrConn,
		Password: ServiceConfig.RedisConfig.Password,
		DB:       ServiceConfig.RedisConfig.DBName,
	})
	status := redisClient.Ping()
	fmt.Println(status)
	//---Redis
}

func main() {
	repo := repoMovie.New(database.Collection("movies"))
	ctrl := ctrlMovie.New(repo)
	handler := handlerMovie.New(ctrl, redisClient)

	router := gin.Default()
	router.POST("/movies", handler.CreateMovie)
	router.GET("/movies", handler.GetMovies)
	router.PUT("/movies/:movieId", handler.UpdateMovie)
	router.DELETE("/movies/:movieId", handler.DeleteMovie)
	router.Run()
}

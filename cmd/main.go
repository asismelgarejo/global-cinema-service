package main

import (
	controllerAuth "api/internal/controller/auth"
	ctrlMovie "api/internal/controller/movie"
	controllerUser "api/internal/controller/user"
	handlerAuth "api/internal/handler/http/auth"
	handlerMovie "api/internal/handler/http/movie"
	handlerUser "api/internal/handler/http/user"
	repoMovie "api/internal/repository/mongodb/movie"
	repositoryUser "api/internal/repository/mongodb/user"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/yaml.v2"

	middleware "api/pkg/middleware"
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

// func init() {
// 	users := map[string]string{
// 		"admin": "fCRmh4Q2J7Rseqkz",
// 		"asis":  "L3nSFRcZzNQ67bcc",
// 	}
// 	ctx := context.Background()
// 	h := sha256.New()
// 	for username, password := range users {
// 		database.Collection("users").InsertOne(ctx, bson.M{
// 			"username": username,
// 			"password": fmt.Sprintf("%x", h.Sum([]byte(password))),
// 			// "password": string(h.Sum([]byte(password))),
// 		})
// 	}
// 	log.Println("Users created!")
// }

func main() {
	router := gin.Default()

	//Movies
	repo := repoMovie.New(database.Collection("movies"))
	ctrl := ctrlMovie.New(repo)
	handler := handlerMovie.New(ctrl, redisClient)

	// User
	repoUser := repositoryUser.New(database.Collection("users"))
	ctrlUser := controllerUser.New(repoUser)
	hUser := handlerUser.New(ctrlUser)

	//Authentication
	ctrlAuth := controllerAuth.New(repoUser)
	hSignIn := handlerAuth.New(ctrlAuth)

	authorized := router.Group("/")

	//Authentication
	router.POST("/signIn", hSignIn.SignIn)
	router.POST("/refresh", hSignIn.RefreshToken)

	//Users
	router.POST("/user", hUser.CreateUser)

	//Movies
	router.GET("/movies", handler.GetMovies)
	router.GET("/movies/:movieId", handler.GetMovie)

	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/movies", handler.CreateMovie)
		authorized.PUT("/movies/:movieId", handler.UpdateMovie)
		authorized.DELETE("/movies/:movieId", handler.DeleteMovie)
	}

	router.Run()
}

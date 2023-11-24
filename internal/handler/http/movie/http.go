package movie

import (
	controller "api/internal/controller/movie"
	"api/pkg/model"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type MovieHandler struct {
	ctrl        *controller.MovieController
	redisClient *redis.Client
}

func New(ctrl *controller.MovieController, redisClient *redis.Client) *MovieHandler {
	return &MovieHandler{ctrl: ctrl, redisClient: redisClient}
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	movieData := model.Movie{}

	if err := c.ShouldBindJSON(&movieData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dataResponse, err := h.ctrl.CreateMovie(c.Request.Context(), movieData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("Remove data from Redis")
	h.redisClient.Del("movies")
	c.JSON(http.StatusOK, gin.H{
		"data": dataResponse,
	})
}
func (h *MovieHandler) GetMovies(c *gin.Context) {

	val, err := h.redisClient.Get("movies").Result()
	if err == redis.Nil {
		log.Println("Request to database!")
		dataResponse, err := h.ctrl.GetMovies(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data, _ := json.Marshal(dataResponse)
		h.redisClient.Set("movies", string(data), 0)
		c.JSON(http.StatusOK, gin.H{
			"data": dataResponse,
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		log.Println("Request to redis")
		var movies []*model.Movie
		json.Unmarshal([]byte(val), &movies)
		c.JSON(http.StatusOK, gin.H{
			"data": movies,
		})
	}
}
func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	movieId := c.Params.ByName("movieId")
	movieData := model.Movie{}
	if err := c.ShouldBindJSON(&movieData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	dataResponse, err := h.ctrl.UpdateMovie(c.Request.Context(), movieId, movieData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.redisClient.Del("movies")
	c.JSON(http.StatusOK, gin.H{
		"data": dataResponse,
	})
}
func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	movieId := c.Params.ByName("movieId")

	dataResponse, err := h.ctrl.DeleteMovie(c.Request.Context(), movieId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("Remove data from Redis")
	h.redisClient.Del("movies")
	c.JSON(http.StatusOK, gin.H{
		"data": dataResponse,
	})
}
func (h *MovieHandler) GetMovie(c *gin.Context) {
	movieId := c.Params.ByName("movieId")

	dataResponse, err := h.ctrl.GetMovie(c.Request.Context(), movieId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": dataResponse,
	})
}

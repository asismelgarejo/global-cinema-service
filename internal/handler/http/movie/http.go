package movie

import (
	controller "api/internal/controller/movie"
	"api/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	ctrl *controller.MovieController
}

func New(ctrl *controller.MovieController) *MovieHandler {
	return &MovieHandler{ctrl: ctrl}
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

	c.JSON(http.StatusOK, gin.H{
		"data": dataResponse,
	})
}
func (h *MovieHandler) GetMovies(c *gin.Context) {
	dataResponse, err := h.ctrl.GetMovies(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dataResponse,
	})
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

	c.JSON(http.StatusOK, gin.H{
		"data": dataResponse,
	})
}

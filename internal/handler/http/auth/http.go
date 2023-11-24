package auth

import (
	"api/internal/controller/auth"
	"api/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerAuth struct {
	ctrl *auth.ControllerAuth
}

func New(ctrl *auth.ControllerAuth) *HandlerAuth {
	return &HandlerAuth{ctrl: ctrl}
}

func (a *HandlerAuth) SignIn(c *gin.Context) {
	var userData model.User
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	response, err := a.ctrl.SignIn(c.Request.Context(), userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
func (a *HandlerAuth) RefreshToken(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	response, err := a.ctrl.RefreshToken(c.Request.Context(), authorizationHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

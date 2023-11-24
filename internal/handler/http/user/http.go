package user

import (
	controllerUser "api/internal/controller/user"
	"api/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerUser struct {
	ctrl *controllerUser.ControllerUser
}

func New(ctrl *controllerUser.ControllerUser) *HandlerUser {
	return &HandlerUser{ctrl: ctrl}
}

func (h *HandlerUser) CreateUser(c *gin.Context) {
	var userData model.User
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	response, err := h.ctrl.CreateUser(c.Request.Context(), userData)
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

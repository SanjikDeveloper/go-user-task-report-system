package delivery

import (
	"net/http"
	"user-service/internal/models"

	"github.com/gin-gonic/gin"
)

type signUpInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	ctx := c.Request.Context()

	var input signUpInput

	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, models.ErrInvalidInput)
		return
	}

	role := input.Role
	if role == "" {
		role = "user"
	}

	user := models.User{
		Username: input.Username,
		Password: input.Password,
		Role:     role,
	}

	id, err := h.service.CreateUser(ctx, user)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) signIn(c *gin.Context) {
	ctx := c.Request.Context()
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, models.ErrInvalidInput)
		return
	}

	token, err := h.service.GenerateToken(ctx, input.Username, input.Password)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

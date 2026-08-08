package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/auth"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type registerResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type AuthHandler struct {
	auth *authservice.Service
}

func NewAuthHandler(auth *authservice.Service) *AuthHandler {
	return &AuthHandler{
		auth: auth,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request registerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid request"},
		)
		return
	}

	created, err := h.auth.Register(
		c.Request.Context(),
		request.Email,
		request.Password,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to register user"},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		registerResponse{
			ID:    created.ID.String(),
			Email: created.Email,
		},
	)
}

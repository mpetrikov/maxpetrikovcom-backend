package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	userservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/user"
)

type UserHandler struct {
	users *userservice.Service
}

func NewUserHandler(
	users *userservice.Service,
) *UserHandler {
	return &UserHandler{
		users: users,
	}
}

type meResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *UserHandler) Me(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "invalid authenticated user"},
		)
		return
	}

	currentUser, err := h.users.FindByID(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to get user"},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		meResponse{
			ID:    currentUser.ID.String(),
			Email: currentUser.Email,
			Role:  string(currentUser.Role),
		},
	)
}

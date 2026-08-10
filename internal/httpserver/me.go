package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) me(c *gin.Context) {
	userID, exists := c.Get(userIDContextKey)
	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "user not found in context"},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"user_id": userID,
		},
	)
}

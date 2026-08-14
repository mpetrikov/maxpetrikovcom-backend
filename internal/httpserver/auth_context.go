package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func authenticatedUserID(
	c *gin.Context,
) (uuid.UUID, bool) {
	userIDValue, exists := c.Get(userIDContextKey)
	if !exists {
		return uuid.Nil, false
	}

	userIDString, ok := userIDValue.(string)
	if !ok {
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, false
	}

	return userID, true
}

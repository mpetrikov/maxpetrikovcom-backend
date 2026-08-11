package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/role"
)

func requireRole(allowedRoles ...role.Name) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get(roleContextKey)
		if !exists {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "role not found"},
			)
			return
		}

		currentRole, ok := value.(role.Name)
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "invalid role"},
			)
			return
		}

		for _, allowedRole := range allowedRoles {
			if currentRole == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{"error": "forbidden"},
		)
	}
}

package httpserver

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	authservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/auth"
)

const (
	userIDContextKey = "user_id"
	roleContextKey   = "role"
)

func authMiddleware(
	tokens *authservice.TokenService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")

		if authorization == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "authorization header is required"},
			)
			return
		}

		parts := strings.SplitN(authorization, " ", 2)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			parts[1] == "" {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid authorization header"},
			)
			return
		}

		claims, err := tokens.ParseAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid access token"},
			)
			return
		}

		c.Set(userIDContextKey, claims.Subject)
		c.Set(roleContextKey, claims.Role)

		c.Next()
	}
}

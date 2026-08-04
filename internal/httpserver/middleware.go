package httpserver

import (
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		startedAt := time.Now()

		context.Next()

		s.logger.Info(
			"http request",
			"method", context.Request.Method,
			"path", context.Request.URL.Path,
			"status", context.Writer.Status(),
			"duration", time.Since(startedAt),
			"client_ip", context.ClientIP(),
		)
	}
}

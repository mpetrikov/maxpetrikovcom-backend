package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type healthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func (s *Server) health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(
		c.Request.Context(),
		2*time.Second,
	)
	defer cancel()

	if err := s.db.Ping(ctx); err != nil {
		c.JSON(
			http.StatusServiceUnavailable,
			healthResponse{
				Status:   "error",
				Database: "unavailable",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		healthResponse{
			Status:   "ok",
			Database: "ok",
		},
	)
}

package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthResponse struct {
	Status string `json:"status"`
}

func (s *Server) health(context *gin.Context) {
	context.JSON(
		http.StatusOK,
		healthResponse{
			Status: "ok",
		},
	)
}

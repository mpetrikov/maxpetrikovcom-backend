package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func adminTest(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "admin access granted",
		},
	)
}

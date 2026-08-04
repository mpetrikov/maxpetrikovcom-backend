package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	logger *slog.Logger
	router *gin.Engine
}

func New(logger *slog.Logger) *Server {
	router := gin.New()

	server := &Server{
		logger: logger,
		router: router,
	}

	server.registerMiddleware()
	server.registerRoutes()

	return server
}

func (s *Server) Handler() http.Handler {
	return s.router
}

func (s *Server) registerMiddleware() {
	s.router.Use(gin.Recovery())
	s.router.Use(s.loggingMiddleware())
}

func (s *Server) registerRoutes() {
	s.router.GET("/health", s.health)
}

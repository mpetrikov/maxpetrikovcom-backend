package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	logger *slog.Logger
	router *gin.Engine
	db     *pgxpool.Pool
}

func New(logger *slog.Logger, db *pgxpool.Pool) *Server {
	router := gin.New()

	server := &Server{
		logger: logger,
		router: router,
		db:     db,
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

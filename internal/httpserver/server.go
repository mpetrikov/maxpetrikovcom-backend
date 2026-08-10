package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	logger      *slog.Logger
	router      *gin.Engine
	db          *pgxpool.Pool
	authHandler *AuthHandler
}

func New(logger *slog.Logger, db *pgxpool.Pool, authHandler *AuthHandler) *Server {
	router := gin.New()

	server := &Server{
		logger:      logger,
		router:      router,
		db:          db,
		authHandler: authHandler,
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

	auth := s.router.Group("/auth")
	{
		auth.POST("/register", s.authHandler.Register)
		auth.GET("/google", s.authHandler.GoogleLogin)
		auth.GET("/google/callback", s.authHandler.GoogleCallback)
	}
}

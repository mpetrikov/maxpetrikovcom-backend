package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/role"
	authservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/auth"
)

type Server struct {
	logger      *slog.Logger
	router      *gin.Engine
	db          *pgxpool.Pool
	authHandler *AuthHandler
	userHandler *UserHandler
	tokens      *authservice.TokenService
}

func New(logger *slog.Logger,
	db *pgxpool.Pool,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	tokens *authservice.TokenService,
) *Server {
	router := gin.New()

	server := &Server{
		logger:      logger,
		router:      router,
		db:          db,
		authHandler: authHandler,
		userHandler: userHandler,
		tokens:      tokens,
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

		auth.POST("/refresh", s.authHandler.Refresh)
		auth.POST("/logout", s.authHandler.Logout)
	}

	protected := s.router.Group("/")
	protected.Use(authMiddleware(s.tokens))
	{
		protected.GET("/me", s.userHandler.Me)

		admin := protected.Group("/admin")
		admin.Use(requireRole(role.Admin))
		{
			admin.GET("/test", adminTest)
		}
	}
}

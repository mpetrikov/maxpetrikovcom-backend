package httpserver

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/user"
	authservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/auth"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type registerResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type AuthHandler struct {
	auth   *authservice.Service
	google *authservice.GoogleOAuth
}

func NewAuthHandler(auth *authservice.Service, google *authservice.GoogleOAuth) *AuthHandler {
	return &AuthHandler{
		auth:   auth,
		google: google,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request registerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid request"},
		)
		return
	}

	created, err := h.auth.Register(
		c.Request.Context(),
		request.Email,
		request.Password,
	)
	if err != nil {
		if errors.Is(err, user.ErrEmailAlreadyExists) {
			c.JSON(
				http.StatusConflict,
				gin.H{"error": "email already registered"},
			)
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to register user"},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		registerResponse{
			ID:    created.ID.String(),
			Email: created.Email,
		},
	)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state, err := authservice.GenerateOAuthState()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to start OAuth flow"},
		)
		return
	}

	c.SetCookie(
		"oauth_state",
		state,
		600,
		"/",
		"",
		false, //local HTTP
		true,  // HttpOnly
	)

	c.Redirect(
		http.StatusTemporaryRedirect,
		h.google.AuthorizationURL(state),
	)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	queryState := c.Query("state")

	cookieState, err := c.Cookie("oauth_state")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "oauth state cookie not found"},
		)
		return
	}

	if queryState == "" || queryState != cookieState {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid oauth state"},
		)
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "oauth code not found"},
		)
		return
	}

	// clear state, it is needed only for one OAuth flow
	c.SetCookie(
		"oauth_state",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	googleUser, err := h.google.GetUser(
		c.Request.Context(),
		code,
	)
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			gin.H{"error": "failed to authenticate with Google"},
		)
		return
	}

	authenticatedUser, err := h.auth.LoginWithGoogle(
		c.Request.Context(),
		googleUser,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to authenticate user"},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"id":    authenticatedUser.ID,
			"email": authenticatedUser.Email,
		},
	)
}

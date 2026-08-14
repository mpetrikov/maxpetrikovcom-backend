package httpserver

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/labsession"

	labsessionservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/labsession"
)

type LabSessionHandler struct {
	labSessions *labsessionservice.Service
}

type createLabSessionRequest struct {
	LabSlug string `json:"lab_slug" binding:"required"`
}

func NewLabSessionHandler(
	labSessions *labsessionservice.Service,
) *LabSessionHandler {
	return &LabSessionHandler{
		labSessions: labSessions,
	}
}

type labSessionResponse struct {
	ID     string `json:"id"`
	LabID  string `json:"lab_id"`
	Status string `json:"status"`

	Namespace *string `json:"namespace,omitempty"`
	PodName   *string `json:"pod_name,omitempty"`

	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	ExpiresAt  time.Time  `json:"expires_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`

	FailureReason *string `json:"failure_reason,omitempty"`
}

func toLabSessionResponse(
	labSession labsession.Session,
) labSessionResponse {
	return labSessionResponse{
		ID:            labSession.ID.String(),
		LabID:         labSession.LabID.String(),
		Status:        string(labSession.Status),
		Namespace:     labSession.Namespace,
		PodName:       labSession.PodName,
		CreatedAt:     labSession.CreatedAt,
		StartedAt:     labSession.StartedAt,
		ExpiresAt:     labSession.ExpiresAt,
		FinishedAt:    labSession.FinishedAt,
		FailureReason: labSession.FailureReason,
	}
}

func (h *LabSessionHandler) Create(
	c *gin.Context,
) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	var request createLabSessionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid request"},
		)
		return
	}

	labSession, err := h.labSessions.Create(
		c.Request.Context(),
		request.LabSlug,
		userID,
	)

	if errors.Is(err, labsession.ErrNotFound) {
		c.JSON(
			http.StatusNotFound,
			gin.H{"error": "lab not found"},
		)
		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to create lab session"},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		toLabSessionResponse(labSession),
	)
}

func (h *LabSessionHandler) Get(
	c *gin.Context,
) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	labSessionID, err := uuid.Parse(
		c.Param("id"),
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid lab session id"},
		)
		return
	}

	labSession, err := h.labSessions.Get(
		c.Request.Context(),
		labSessionID,
		userID,
	)

	if errors.Is(err, labsession.ErrNotFound) {
		c.JSON(
			http.StatusNotFound,
			gin.H{"error": "lab session not found"},
		)
		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to get lab session"},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		toLabSessionResponse(labSession),
	)
}

func (h *LabSessionHandler) List(
	c *gin.Context,
) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	labSessions, err := h.labSessions.ListForUser(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to get lab sessions"},
		)
		return
	}

	response := make(
		[]labSessionResponse,
		0,
		len(labSessions),
	)

	for _, labSession := range labSessions {
		response = append(
			response,
			toLabSessionResponse(labSession),
		)
	}

	c.JSON(http.StatusOK, response)
}

func (h *LabSessionHandler) Stop(
	c *gin.Context,
) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "unauthorized"},
		)
		return
	}

	labSessionID, err := uuid.Parse(
		c.Param("id"),
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid lab session id"},
		)
		return
	}

	err = h.labSessions.Stop(
		c.Request.Context(),
		labSessionID,
		userID,
	)

	if errors.Is(err, labsession.ErrNotFound) {
		c.JSON(
			http.StatusNotFound,
			gin.H{"error": "lab session not found"},
		)
		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to stop lab session"},
		)
		return
	}

	c.Status(http.StatusNoContent)
}

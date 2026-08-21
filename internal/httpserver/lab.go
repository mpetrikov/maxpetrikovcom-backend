package httpserver

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainlab "github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/lab"
	labservice "github.com/maxpetrikov/maxpetrikovcom-backend/internal/service/lab"
)

type LabHandler struct {
	labs *labservice.Service
}

func NewLabHandler(
	labs *labservice.Service,
) *LabHandler {
	return &LabHandler{
		labs: labs,
	}
}

type createLabRequest struct {
	Slug           string               `json:"slug" binding:"required"`
	Title          string               `json:"title" binding:"required"`
	Description    string               `json:"description" binding:"required"`
	Difficulty     domainlab.Difficulty `json:"difficulty" binding:"required,oneof=easy medium hard"`
	TimeoutMinutes int                  `json:"timeout_minutes" binding:"required,min=1,max=240"`
	Image          string               `json:"image" binding:"required"`
	CPULimit       string               `json:"cpu_limit" binding:"required"`
	MemoryLimit    string               `json:"memory_limit" binding:"required"`
	IsPublished    bool                 `json:"is_published"`
}

type labResponse struct {
	ID             string               `json:"id"`
	Slug           string               `json:"slug"`
	Title          string               `json:"title"`
	Description    string               `json:"description"`
	Difficulty     domainlab.Difficulty `json:"difficulty"`
	TimeoutMinutes int                  `json:"timeout_minutes"`
	Image          string               `json:"image"`
	CPULimit       string               `json:"cpu_limit"`
	MemoryLimit    string               `json:"memory_limit"`
	IsPublished    bool                 `json:"is_published"`
}

func toLabResponse(item domainlab.Lab) labResponse {
	return labResponse{
		ID:             item.ID.String(),
		Slug:           item.Slug,
		Title:          item.Title,
		Description:    item.Description,
		Difficulty:     item.Difficulty,
		TimeoutMinutes: item.TimeoutMinutes,
		Image:          item.Image,
		CPULimit:       item.CPULimit,
		MemoryLimit:    item.MemoryLimit,
		IsPublished:    item.IsPublished,
	}
}

func (h *LabHandler) Create(c *gin.Context) {
	var request createLabRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid request"},
		)
		return
	}

	created, err := h.labs.Create(
		c.Request.Context(),
		labservice.CreateInput{
			Slug:           request.Slug,
			Title:          request.Title,
			Description:    request.Description,
			Difficulty:     request.Difficulty,
			TimeoutMinutes: request.TimeoutMinutes,
			Image:          request.Image,
			CPULimit:       request.CPULimit,
			MemoryLimit:    request.MemoryLimit,
			IsPublished:    request.IsPublished,
		},
	)

	if errors.Is(err, domainlab.ErrSlugAlreadyExists) {
		c.JSON(
			http.StatusConflict,
			gin.H{"error": "lab slug already exists"},
		)
		return
	}

	if errors.Is(err, domainlab.ErrInvalidInput) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid lab runtime config"},
		)
		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to create lab"},
		)
		return
	}

	c.JSON(http.StatusCreated, toLabResponse(created))
}

func (h *LabHandler) List(c *gin.Context) {
	items, err := h.labs.ListPublished(c.Request.Context())
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to get labs"},
		)
		return
	}

	response := make([]labResponse, 0, len(items))

	for _, item := range items {
		response = append(response, toLabResponse(item))
	}

	c.JSON(http.StatusOK, response)
}

func (h *LabHandler) Get(c *gin.Context) {
	item, err := h.labs.FindBySlug(
		c.Request.Context(),
		c.Param("slug"),
	)

	if errors.Is(err, domainlab.ErrNotFound) ||
		(err == nil && !item.IsPublished) {

		c.JSON(
			http.StatusNotFound,
			gin.H{"error": "lab not found"},
		)
		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to get lab"},
		)
		return
	}

	c.JSON(http.StatusOK, toLabResponse(item))
}

package controllers

import (
	"backend/models"
	"backend/repositories"
	"backend/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/microcosm-cc/bluemonday"
)

type Drawing struct {
	repo      repositories.DrawingRepository
	service   *services.Drawing
	ugcPolicy *bluemonday.Policy
}

func NewDrawing(repo repositories.DrawingRepository, service *services.Drawing) *Drawing {
	return &Drawing{
		repo:      repo,
		service:   service,
		ugcPolicy: bluemonday.UGCPolicy(),
	}
}

func (ctrl *Drawing) GetDrawings(c *gin.Context) {
	projectIDStr := c.Query("project_id")
	projectID := 0
	if projectIDStr != "" {
		projectID, _ = strconv.Atoi(projectIDStr)
	}

	drawings, err := ctrl.repo.GetByProject(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch drawings"})
		return
	}
	c.JSON(http.StatusOK, drawings)
}

type CreateDrawingRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"max=500"`
	ProjectID   uint   `json:"project_id" binding:"required"`
}

func (ctrl *Drawing) CreateDrawing(c *gin.Context) {
	var req CreateDrawingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErrorMessage(err)})
		return
	}

	sanitizedTitle := ctrl.ugcPolicy.Sanitize(req.Title)
	sanitizedDesc := ctrl.ugcPolicy.Sanitize(req.Description)

	drawing := models.Drawing{
		Title:        sanitizedTitle,
		Description:  sanitizedDesc,
		ProjectID:    req.ProjectID,
		CurrentStage: models.StageUnassigned,
		Version:      0,
		AuthorID:     c.MustGet("user_id").(uint),
	}

	if err := ctrl.repo.Create(&drawing); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "A drawing with this title already exists in this project"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create drawing"})
		return
	}

	c.JSON(http.StatusCreated, drawing)
}

func (ctrl *Drawing) ClaimDrawing(c *gin.Context) {
	ctrl.handleWorkflowAction(c, models.ActionClaim)
}

func (ctrl *Drawing) SubmitDrawing(c *gin.Context) {
	ctrl.handleWorkflowAction(c, models.ActionSubmit)
}

func (ctrl *Drawing) ReleaseDrawing(c *gin.Context) {
	ctrl.handleWorkflowAction(c, models.ActionRelease)
}

func (ctrl *Drawing) RejectDrawing(c *gin.Context) {
	ctrl.handleWorkflowAction(c, models.ActionReject)
}

func (ctrl *Drawing) handleWorkflowAction(c *gin.Context, action models.Action) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	userID := c.MustGet("user_id").(uint)
	userRole := c.MustGet("role").(string)

	drawing, err := ctrl.service.ProcessWorkflowAction(uint(id), userID, userRole, action)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Action processed successfully",
		"drawing": drawing,
	})
}

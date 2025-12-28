package repositories

import (
	"backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DrawingRepository interface
type DrawingRepository interface {
	Get(id uint) (*models.Drawing, error)
	GetForUpdate(id uint) (*models.Drawing, error)
	Update(drawing *models.Drawing, updates map[string]interface{}) error
	CreateWorkflowLog(log models.WorkflowLog) error
	GetByProject(projectID uint) ([]models.Drawing, error)
	Create(drawing *models.Drawing) error

	// Transaction support
	RunTransaction(fn func(repo DrawingRepository) error) error
}

// GormDrawingRepository implementation
type GormDrawingRepository struct {
	db *gorm.DB
}

func NewDrawingRepository(db *gorm.DB) *GormDrawingRepository {
	return &GormDrawingRepository{db: db}
}

func (r *GormDrawingRepository) Get(id uint) (*models.Drawing, error) {
	var drawing models.Drawing
	if err := r.db.Where("id = ?", id).First(&drawing).Error; err != nil {
		return nil, err
	}
	return &drawing, nil
}

func (r *GormDrawingRepository) GetForUpdate(id uint) (*models.Drawing, error) {
	var drawing models.Drawing
	if err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&drawing).Error; err != nil {
		return nil, err
	}
	return &drawing, nil
}

func (r *GormDrawingRepository) Update(drawing *models.Drawing, updates map[string]interface{}) error {
	result := r.db.Model(drawing).
		Where("id = ? AND version = ?", drawing.ID, drawing.Version).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *GormDrawingRepository) CreateWorkflowLog(log models.WorkflowLog) error {
	return r.db.Create(&log).Error
}

func (r *GormDrawingRepository) GetByProject(projectID uint) ([]models.Drawing, error) {
	var drawings []models.Drawing
	query := r.db.Preload("Assignee")
	if projectID != 0 {
		query = query.Where("project_id = ?", projectID)
	}
	err := query.Find(&drawings).Error
	return drawings, err
}

func (r *GormDrawingRepository) Create(drawing *models.Drawing) error {
	return r.db.Create(drawing).Error
}

func (r *GormDrawingRepository) RunTransaction(fn func(repo DrawingRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := NewDrawingRepository(tx)
		return fn(txRepo)
	})
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleDrafter   UserRole = "drafter"
	RoleShiftLead UserRole = "shift_lead"
	RoleFinalQC   UserRole = "final_qc"
)

type Stage string

const (
	StageUnassigned Stage = "unassigned"
	StageDrafting   Stage = "drafting"
	StageFirstQC    Stage = "first_qc"
	StageFinalQC    Stage = "final_qc"
	StageApproved   Stage = "approved"
)

type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name" binding:"required"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Members  []User    `gorm:"many2many:project_members;" json:"members,omitempty"`
	Drawings []Drawing `json:"drawings,omitempty"`
}

type ProjectMember struct {
	ProjectID uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"primaryKey"`
	Role      UserRole  `gorm:"not null"` // Role specifically within this project
	JoinedAt  time.Time `gorm:"autoCreateTime"`
}

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username" binding:"required"`
	PasswordHash string         `json:"-"`
	Role         UserRole       `gorm:"not null" json:"role"` // System-wide default role
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Drawing struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Title       string  `gorm:"uniqueIndex:idx_project_title;not null" json:"title" binding:"required"`
	Description string  `json:"description"`
	ProjectID   uint    `gorm:"uniqueIndex:idx_project_title;not null" json:"project_id"`
	Project     Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`

	AuthorID uint `gorm:"index" json:"author_id"` // Creator of the drawing (Nullable for existing)
	Author   User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`

	CurrentStage Stage `gorm:"index;not null;default:'unassigned'" json:"current_stage"`
	AssigneeID   *uint `gorm:"index" json:"assignee_id"`
	Assignee     *User `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`

	Revision   int    `gorm:"not null;default:1" json:"revision"` // Business Revision (increases on rework/submit)
	Version    int64  `gorm:"not null;default:0" json:"version"`  // Technical Concurrency Lock
	DrawingURL string `json:"drawing_url"`                        // Link to S3/CDN file

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type WorkflowLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DrawingID uint      `gorm:"not null;index" json:"drawing_id"`
	ActorID   uint      `gorm:"not null" json:"actor_id"`
	Actor     User      `gorm:"foreignKey:ActorID" json:"actor"`
	Action    string    `json:"action"` // e.g., "claimed", "submitted", "rejected"
	FromStage Stage     `json:"from_stage"`
	ToStage   Stage     `json:"to_stage"`
	Comment   string    `json:"comment"`
	Timestamp time.Time `gorm:"autoCreateTime" json:"timestamp"`
}

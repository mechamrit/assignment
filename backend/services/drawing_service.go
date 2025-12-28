package services

import (
	"backend/models"
	"backend/repositories"
	"fmt"
)

// Auditor interface abstraction
type Auditor interface {
	ProduceAuditLog(workflowLog models.WorkflowLog)
}

// Broadcaster interface abstraction
type Broadcaster interface {
	BroadcastEvent(projectID uint, eventType string, payload interface{})
}

type Drawing struct {
	repo        repositories.DrawingRepository
	auditor     Auditor
	broadcaster Broadcaster
}

func NewDrawing(repo repositories.DrawingRepository, auditor Auditor, broadcaster Broadcaster) *Drawing {
	return &Drawing{
		repo:        repo,
		auditor:     auditor,
		broadcaster: broadcaster,
	}
}

func (s *Drawing) ProcessWorkflowAction(id uint, userID uint, userRole string, action models.Action) (*models.Drawing, error) {
	var drawing models.Drawing
	var workflowLog models.WorkflowLog

	err := s.repo.RunTransaction(func(txRepo repositories.DrawingRepository) error {
		d, err := txRepo.GetForUpdate(id)
		if err != nil {
			return err
		}
		drawing = *d

		// Validation for Claim/Submit/Release/Reject
		if action != models.ActionClaim {
			// Admins can perform any action without being the assignee
			if userRole != string(models.RoleAdmin) {
				if drawing.AssigneeID == nil || *drawing.AssigneeID != userID {
					return fmt.Errorf("drawing not assigned to user or unassigned")
				}
			}
		}

		nextStage, err := models.GetNextState(drawing.CurrentStage, action, models.UserRole(userRole))
		if err != nil {
			return err
		}

		updates := map[string]interface{}{
			"current_stage": nextStage,
			"version":       drawing.Version + 1,
		}

		// Handle Assignee and Revision logic based on action
		if action == models.ActionClaim {
			if drawing.AssigneeID != nil {
				return fmt.Errorf("drawing already claimed")
			}
			updates["assignee_id"] = userID
		} else if action == models.ActionSubmit || action == models.ActionReject {
			// Submit or Reject increments the business revision
			updates["assignee_id"] = nil
			updates["revision"] = drawing.Revision + 1
		} else {
			// Release only clears assignee
			updates["assignee_id"] = nil
		}

		if err := txRepo.Update(&drawing, updates); err != nil {
			if err.Error() == "record not found" {
				return fmt.Errorf("concurrent update detected")
			}
			return err
		}

		workflowLog = models.WorkflowLog{
			DrawingID: drawing.ID,
			ActorID:   userID,
			Action:    string(action),
			FromStage: drawing.CurrentStage,
			ToStage:   nextStage,
		}
		return txRepo.CreateWorkflowLog(workflowLog)
	})

	if err != nil {
		return nil, err
	}

	// Post-transaction tasks (Async)
	go func() {
		s.auditor.ProduceAuditLog(workflowLog)
		s.broadcaster.BroadcastEvent(drawing.ProjectID, fmt.Sprintf("DRAWING_%s", action), drawing)
	}()

	return &drawing, nil
}

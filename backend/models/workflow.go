package models

import (
	"errors"
)

var (
	ErrInvalidTransition = errors.New("invalid state transition")
	ErrUnauthorizedRole  = errors.New("role not authorized for this action")
)

type Action string

const (
	ActionClaim   Action = "claim"
	ActionSubmit  Action = "submit"
	ActionRelease Action = "release"
	ActionReject  Action = "reject"
)

// Workflow defines the transition rules
type transition struct {
	From   Stage
	Action Action
	To     Stage
	Role   UserRole
}

var workflowRules = []transition{
	// Admin flow
	{From: StageUnassigned, Action: ActionClaim, To: StageUnassigned, Role: RoleAdmin},
	{From: StageUnassigned, Action: ActionSubmit, To: StageDrafting, Role: RoleAdmin},
	{From: StageUnassigned, Action: ActionRelease, To: StageUnassigned, Role: RoleAdmin},
	{From: StageDrafting, Action: ActionRelease, To: StageDrafting, Role: RoleAdmin},
	{From: StageFirstQC, Action: ActionRelease, To: StageFirstQC, Role: RoleAdmin},
	{From: StageFinalQC, Action: ActionRelease, To: StageFinalQC, Role: RoleAdmin},
	{From: StageFirstQC, Action: ActionReject, To: StageDrafting, Role: RoleAdmin},
	{From: StageFinalQC, Action: ActionReject, To: StageDrafting, Role: RoleAdmin},

	// Drafting flow
	{From: StageUnassigned, Action: ActionClaim, To: StageDrafting, Role: RoleDrafter},
	{From: StageDrafting, Action: ActionClaim, To: StageDrafting, Role: RoleDrafter},
	{From: StageDrafting, Action: ActionSubmit, To: StageFirstQC, Role: RoleDrafter},
	{From: StageDrafting, Action: ActionRelease, To: StageDrafting, Role: RoleDrafter},

	// First QC flow
	{From: StageFirstQC, Action: ActionClaim, To: StageFirstQC, Role: RoleShiftLead},
	{From: StageFirstQC, Action: ActionSubmit, To: StageFinalQC, Role: RoleShiftLead},
	{From: StageFirstQC, Action: ActionRelease, To: StageFirstQC, Role: RoleShiftLead},
	{From: StageFirstQC, Action: ActionReject, To: StageDrafting, Role: RoleShiftLead},

	// Final QC flow
	{From: StageFinalQC, Action: ActionClaim, To: StageFinalQC, Role: RoleFinalQC},
	{From: StageFinalQC, Action: ActionSubmit, To: StageApproved, Role: RoleFinalQC},
	{From: StageFinalQC, Action: ActionRelease, To: StageFinalQC, Role: RoleFinalQC},
	{From: StageFinalQC, Action: ActionReject, To: StageDrafting, Role: RoleFinalQC},
}

// GetNextState validates the transition and returns the next state
func GetNextState(current Stage, action Action, role UserRole) (Stage, error) {
	foundAction := false
	for _, t := range workflowRules {
		if t.From == current && t.Action == action {
			foundAction = true
			if t.Role == role {
				return t.To, nil
			}
		}
	}
	if foundAction {
		return "", ErrUnauthorizedRole
	}
	return "", ErrInvalidTransition
}

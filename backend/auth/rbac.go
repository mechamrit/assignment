package auth

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var Enforcer *casbin.Enforcer

func InitCasbin(db *gorm.DB) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin adapter: %v", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	modelPath := filepath.Join(filepath.Dir(filename), "rbac_model.conf")

	Enforcer, err = casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		log.Fatalf("Failed to create Casbin enforcer: %v", err)
	}

	err = Enforcer.LoadPolicy()
	if err != nil {
		log.Fatalf("Failed to load Casbin policy: %v", err)
	}

	setupDefaultPolicies()
}

func setupDefaultPolicies() {
	// Roles: admin, drafter, shift_lead, final_qc
	// Actions: create, view, claim, submit, approve, release
	// Resources: drawings

	// Admin can do everything
	Enforcer.AddNamedPolicy("p", "admin", "drawings", "*")

	// Drafter
	Enforcer.AddNamedPolicy("p", "drafter", "drawings", "view")
	Enforcer.AddNamedPolicy("p", "drafter", "drawings", "claim")
	Enforcer.AddNamedPolicy("p", "drafter", "drawings", "submit")
	Enforcer.AddNamedPolicy("p", "drafter", "drawings", "release")

	// Shift Lead
	Enforcer.AddNamedPolicy("p", "shift_lead", "drawings", "view")
	Enforcer.AddNamedPolicy("p", "shift_lead", "drawings", "claim")
	Enforcer.AddNamedPolicy("p", "shift_lead", "drawings", "submit")
	Enforcer.AddNamedPolicy("p", "shift_lead", "drawings", "release")
	Enforcer.AddNamedPolicy("p", "shift_lead", "drawings", "reject")

	// Final QC
	Enforcer.AddNamedPolicy("p", "final_qc", "drawings", "view")
	Enforcer.AddNamedPolicy("p", "final_qc", "drawings", "claim")
	Enforcer.AddNamedPolicy("p", "final_qc", "drawings", "submit")
	Enforcer.AddNamedPolicy("p", "final_qc", "drawings", "approve")
	Enforcer.AddNamedPolicy("p", "final_qc", "drawings", "release")
	Enforcer.AddNamedPolicy("p", "final_qc", "drawings", "reject")

	err := Enforcer.SavePolicy()
	if err != nil {
		log.Printf("Failed to save Casbin policy: %v", err)
	}
}

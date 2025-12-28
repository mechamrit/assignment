package main

import (
	"backend/config"
	"backend/database"
	"backend/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.LoadConfig()
	database.InitDB(cfg.DBURL)

	// Seed Projects
	p1 := models.Project{Name: "Aerodynamics Suite", Description: "Advanced CFD and wind tunnel models."}
	database.DB.Where("name = ?", p1.Name).FirstOrCreate(&p1)

	p2 := models.Project{Name: "Power Unit", Description: "MGU-H and Internal Combustion engine design."}
	database.DB.Where("name = ?", p2.Name).FirstOrCreate(&p2)

	users := []models.User{
		{Username: "admin", Role: models.RoleAdmin},
		{Username: "drafter1", Role: models.RoleDrafter},
		{Username: "drafter2", Role: models.RoleDrafter},
		{Username: "lead1", Role: models.RoleShiftLead},
		{Username: "lead2", Role: models.RoleShiftLead},
		{Username: "qc1", Role: models.RoleFinalQC},
	}

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	for i := range users {
		users[i].PasswordHash = string(hashedPassword)
		if err := database.DB.Where("username = ?", users[i].Username).FirstOrCreate(&users[i]).Error; err != nil {
			log.Printf("Error seeding user %s: %v", users[i].Username, err)
		}

		// Assign everyone to Project 1 for testing
		database.DB.Where(models.ProjectMember{ProjectID: p1.ID, UserID: users[i].ID}).FirstOrCreate(&models.ProjectMember{
			ProjectID: p1.ID,
			UserID:    users[i].ID,
			Role:      users[i].Role,
		})
	}

	drawings := []models.Drawing{
		{Title: "Front Wing v1", ProjectID: p1.ID, CurrentStage: models.StageUnassigned, Revision: 1, AuthorID: users[0].ID},
		{Title: "Rear Diffuser Main", ProjectID: p1.ID, CurrentStage: models.StageUnassigned, Revision: 1, AuthorID: users[0].ID},
		{Title: "Engine Block Block", ProjectID: p2.ID, CurrentStage: models.StageUnassigned, Revision: 1, AuthorID: users[0].ID},
		{Title: "Turbocharger Housing", ProjectID: p2.ID, CurrentStage: models.StageDrafting, Revision: 2, AuthorID: users[0].ID},
	}

	for i := range drawings {
		if err := database.DB.Where("title = ? AND project_id = ?", drawings[i].Title, drawings[i].ProjectID).FirstOrCreate(&drawings[i]).Error; err != nil {
			log.Printf("Error seeding drawing %s: %v", drawings[i].Title, err)
		}
	}

	log.Println("Database seeded successfully with projects and members")
}

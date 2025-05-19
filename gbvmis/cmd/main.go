package main

import (
	"gbvmis/config"
	"gbvmis/internal/api"
	"gbvmis/internal/handlers"
	"gbvmis/internal/models"
	"gbvmis/internal/repository"
	"gbvmis/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate all models
	db.AutoMigrate(
		&models.Victim{},
		&models.Accused{},
		&models.MedicalExamination{},
		&models.PoliceOfficer{},
		&models.PoliceStation{},
	)

	// Repository layer
	victimRepo := &repository.GormVictimRepository{DB: db}
	accusedRepo := &repository.GormAccusedRepository{DB: db}
	examRepo := &repository.GormMedicalExaminationRepository{DB: db}
	officerRepo := &repository.GormPoliceOfficerRepository{DB: db}
	stationRepo := &repository.GormPoliceStationRepository{DB: db}

	// Service layer
	victimService := &service.VictimServiceImpl{Repo: victimRepo}
	accusedService := &service.AccusedServiceImpl{Repo: accusedRepo}
	examService := &service.MedicalExaminationServiceImpl{Repo: examRepo}
	officerService := &service.PoliceOfficerServiceImpl{Repo: officerRepo}
	stationService := &service.PoliceStationServiceImpl{Repo: stationRepo}

	// Handlers
	victimHandler := &handlers.VictimHandler{Service: victimService}
	accusedHandler := &handlers.AccusedHandler{Service: accusedService}
	examHandler := &handlers.ExaminationHandler{Service: examService}
	officerHandler := &handlers.PoliceOfficerHandler{Service: officerService}
	stationHandler := &handlers.PoliceStationHandler{Service: stationService}

	app := fiber.New()

	// Register all API routes with real handlers
	api.RegisterRoutesWithHandlers(app, victimHandler, accusedHandler, examHandler, officerHandler, stationHandler)

	log.Fatal(app.Listen(":8080"))
}

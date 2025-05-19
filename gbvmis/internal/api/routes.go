package api

import (
	"gbvmis/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// Deprecated: use RegisterRoutesWithHandlers for dependency injection
}

func RegisterRoutesWithHandlers(app *fiber.App,
	victimHandler *handlers.VictimHandler,
	accusedHandler *handlers.AccusedHandler,
	examHandler *handlers.ExaminationHandler,
	officerHandler *handlers.PoliceOfficerHandler,
	stationHandler *handlers.PoliceStationHandler,
) {
	api := app.Group("/api")

	victims := api.Group("/victims")
	victims.Post("/", victimHandler.RegisterVictim)
	victims.Get("/:id", victimHandler.GetVictimByID)
	victims.Get("/", victimHandler.ListVictims)
	victims.Put("/:id", victimHandler.UpdateVictim)
	victims.Delete("/:id", victimHandler.DeleteVictim)

	accused := api.Group("/accused")
	accused.Post("/", accusedHandler.RegisterAccused)
	accused.Get("/:id", accusedHandler.GetAccusedByID)
	accused.Get("/", accusedHandler.ListAccused)
	accused.Put("/:id", accusedHandler.UpdateAccused)
	accused.Delete("/:id", accusedHandler.DeleteAccused)

	examinations := api.Group("/examinations")
	examinations.Post("/", examHandler.RegisterExamination)
	examinations.Get("/:id", examHandler.GetExaminationByID)
	examinations.Get("/", examHandler.ListExaminations)
	examinations.Put("/:id", examHandler.UpdateExamination)
	examinations.Delete("/:id", examHandler.DeleteExamination)

	officers := api.Group("/officers")
	officers.Post("/", officerHandler.RegisterOfficer)
	officers.Get("/:id", officerHandler.GetOfficerByID)
	officers.Get("/", officerHandler.ListOfficers)
	officers.Put("/:id", officerHandler.UpdateOfficer)
	officers.Delete("/:id", officerHandler.DeleteOfficer)

	stations := api.Group("/stations")
	stations.Post("/", stationHandler.RegisterStation)
	stations.Get("/:id", stationHandler.GetStationByID)
	stations.Get("/", stationHandler.ListStations)
	stations.Put("/:id", stationHandler.UpdateStation)
	stations.Delete("/:id", stationHandler.DeleteStation)
}

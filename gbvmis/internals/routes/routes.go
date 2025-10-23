package routes

import (
	"gbvmis/internals/controllers"
	"gbvmis/internals/middleware"
	"gbvmis/internals/repository"
	"gbvmis/internals/service"
	"gbvmis/internals/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes initializes all routes with their respective controllers
func SetupRoute(app *fiber.App, db *gorm.DB) {

	authGroup := app.Group("/api")
	authGroup.Post("/login", controllers.Login(db))
	authGroup.Post("/refresh-token", controllers.RefreshToken)

	// Protected routes
	protected := app.Group("/api", middleware.JWTProtected())
	protected.Get("/me", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*utils.Claims)
		return c.JSON(user)
	})

	victimService := repository.VictimDbService(db)
	victimController := controllers.NewVictimController(victimService)
	protected.Get("/victims", victimController.GetAllVictims)
	protected.Get("/victims/search", victimController.SearchVictims)
	victim := protected.Group("/victim")
	victim.Post("/", victimController.CreateVictim)
	victim.Get("/:id", victimController.GetSingleVictim)
	victim.Put("/:id", victimController.UpdateVictim)
	victim.Delete("/:id", victimController.DeleteVictimByID)

	caseService := repository.CaseDbService(db)
	caseController := controllers.NewCaseController(caseService)
	protected.Get("/cases", caseController.GetAllCases)
	protected.Get("/cases/search", caseController.SearchCases)
	casee := protected.Group("/case")
	casee.Post("/", caseController.CreateCase)
	casee.Get("/:id", caseController.GetSingleCase)
	casee.Put("/:id", caseController.UpdateCase)
	casee.Delete("/:id", caseController.DeleteCaseByID)

	chargeService := repository.ChargeDbService(db)
	chargeController := controllers.NewChargeController(chargeService)
	protected.Get("/charges", chargeController.GetAllCharges)
	protected.Get("/charges/search", chargeController.SearchCharges)
	charge := protected.Group("/charge")
	charge.Post("/", chargeController.CreateCharge)
	charge.Get("/:id", chargeController.GetSingleCharge)
	charge.Put("/:id", chargeController.UpdateCharge)
	charge.Delete("/:id", chargeController.DeleteChargeByID)

	suspectService := repository.SuspectDbService(db)
	suspectController := controllers.NewSuspectController(suspectService)
	protected.Get("/suspects", suspectController.GetAllSuspects)
	protected.Get("/suspects/search", suspectController.SearchSuspects)
	suspect := protected.Group("/suspect")
	suspect.Post("/", suspectController.CreateSuspect)
	suspect.Get("/:id", suspectController.GetSingleSuspect)
	suspect.Put("/:id", suspectController.UpdateSuspect)
	suspect.Delete("/:id", suspectController.DeleteSuspectByID)

	policePostService := repository.PolicePostDbService(db)
	policePostController := controllers.NewPolicePostController(policePostService)
	protected.Get("/police-posts", policePostController.GetAllPolicePosts)
	protected.Get("/police-posts/search", policePostController.SearchPolicePosts)
	policePost := protected.Group("/police-post")
	policePost.Post("/", policePostController.CreatePolicePost)
	policePost.Get("/:id", policePostController.GetSinglePolicePost)
	policePost.Put("/:id", policePostController.UpdatePolicePost)
	policePost.Delete("/:id", policePostController.DeletePolicePostByID)

	policeOfficerService := repository.PoliceOfficerDbService(db)
	policeOfficerController := controllers.NewPoliceOfficerController(policeOfficerService)
	protected.Get("/police-officers", policeOfficerController.GetAllPoliceOfficers)
	protected.Get("/police-officers/search", policeOfficerController.SearchPoliceOfficers)
	policeOfficer := protected.Group("/police-officer")
	policeOfficer.Post("/", policeOfficerController.CreatePoliceOfficer)
	policeOfficer.Get("/:id", policeOfficerController.GetSinglePoliceOfficer)
	policeOfficer.Put("/:id", policeOfficerController.UpdatePoliceOfficer)
	policeOfficer.Delete("/:id", policeOfficerController.DeletePoliceOfficerByID)

	healthFacilityService := repository.HealthFacilityDbService(db)
	healthFacilityController := controllers.NewHealthFacilityController(healthFacilityService)
	protected.Get("/health-facilities", healthFacilityController.GetAllHealthFacilities)
	protected.Get("/health-facilities/search", healthFacilityController.SearchHealthFacilities)
	healthFacility := protected.Group("/health-facility")
	healthFacility.Post("/", healthFacilityController.CreateHealthFacility)
	healthFacility.Get("/:id", healthFacilityController.GetSingleHealthFacility)
	healthFacility.Put("/:id", healthFacilityController.UpdateHealthFacility)
	healthFacility.Delete("/:id", healthFacilityController.DeleteHealthFacilityByID)

	healthPractitionerService := repository.HealthPractitionerDbService(db)
	healthPractitionerController := controllers.NewHealthPractitionerController(healthPractitionerService)
	protected.Get("/health-practitioners", healthPractitionerController.GetAllHealthPractitioners)
	protected.Get("/health-practitioners/search", healthPractitionerController.SearchHealthPractitioners)
	healthPractitioner := protected.Group("/health-practitioner")
	healthPractitioner.Post("/", healthPractitionerController.CreateHealthPractitioner)
	healthPractitioner.Get("/:id", healthPractitionerController.GetSingleHealthPractitioner)
	healthPractitioner.Put("/:id", healthPractitionerController.UpdateHealthPractitioner)
	healthPractitioner.Delete("/:id", healthPractitionerController.DeleteHealthPractitionerByID)

	examinationService := repository.ExaminationDbService(db)
	examinationController := controllers.NewExaminationController(examinationService)
	protected.Get("/examinations", examinationController.GetAllExaminations)
	protected.Get("/examinations/search", examinationController.SearchExaminations)
	examination := protected.Group("/examination")
	examination.Post("/", examinationController.CreateExamination)
	examination.Get("/:id", examinationController.GetSingleExamination)
	examination.Put("/:id", examinationController.UpdateExamination)
	examination.Delete("/:id", examinationController.DeleteExaminationByID)

	policeRolesService := repository.RoleDbService(db)
	policeRolesController := controllers.NewPoliceRolesController(policeRolesService)
	protected.Get("/police-roles", policeRolesController.GetAllPoliceRoles)
	protected.Get("/police-roles/search", policeRolesController.SearchPoliceRoles)
	policeRoles := protected.Group("/police-role")
	policeRoles.Post("/", policeRolesController.CreatePoliceRoles)
	policeRoles.Get("/:id", policeRolesController.GetSinglePoliceRole)
	policeRoles.Put("/:id", policeRolesController.UpdatePoliceRole)
	policeRoles.Delete("/:id", policeRolesController.DeletePoliceRole)

	toxicologyRepo := repository.NewToxicologyReportRepository(db)
	toxicologyService := service.NewToxicologyService(toxicologyRepo)
	toxicologyController := controllers.NewToxicologyForensicController(toxicologyService)
	protected.Get("/toxicology-reports", toxicologyController.GetAll)
	protected.Get("/toxicology-reports-pag", toxicologyController.GetPaginatedReports)
	toxicology := protected.Group("/toxicology-report")
	toxicology.Post("/", toxicologyController.Create)
	toxicology.Get("/:id", toxicologyController.GetByID)
	toxicology.Put("/:id", toxicologyController.Update)
	toxicology.Delete("/:id", toxicologyController.Delete)

	NotFoundRoute(app)
}

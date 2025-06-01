package main

import (
	"gbvmis/internals/database"
	"gbvmis/internals/routes"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "gbvmis/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/swagger" // swagger handler
)

// @title SGBV Police API
// @version 1.0
// @description This is a gbvmis api
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080/api
// @BasePath /
func main() {
	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // 20 MB
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	// Initialize and connect to the database
	db := database.NewDatabase()
	db.Connect()
	defer db.Close() // Ensure database connection is closed when the app shuts down

	// Run migrations and seed data
	db.Migrate()
	db.Seed()

	var sessionStore = session.New()
	app.Use(func(c *fiber.Ctx) error {
		sess, err := sessionStore.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Session error")
		}
		c.Locals("session", sess)
		return c.Next()
	})
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Setup routes
	routes.SetupRoute(app, db.GetDB())

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Error during shutdown: %v", err)
		}
	}()

	// Start the server
	app.Listen(":8080")
}

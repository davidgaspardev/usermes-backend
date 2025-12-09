package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	resourceUseCase "github.com/davidgaspardev/usermes-backend/internal/modules/production/resource/application/usecase"
	resourceHttp "github.com/davidgaspardev/usermes-backend/internal/modules/production/resource/infrastructure/adapter/input/http"
	resourcePersistence "github.com/davidgaspardev/usermes-backend/internal/modules/production/resource/infrastructure/adapter/output/persistence"
)

// This is an example of how to integrate the Resource module into your application
func main() {
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "UserMes API - Resource Module Example",
		ErrorHandler: customErrorHandler,
	})

	// Add middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Initialize Resource module
	// 1. Create repository (in-memory for this example)
	resourceRepo := resourcePersistence.NewMemoryResourceRepository()

	// 2. Create service with business logic
	resourceService := resourceUseCase.NewResourceService(resourceRepo)

	// 3. Setup HTTP routes
	resourceRoutes := resourceHttp.NewResourceRoutes(resourceService)
	resourceRoutes.SetupRoutes(app)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"module": "resource",
		})
	})

	// API documentation endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "UserMes API - Resource Module",
			"version": "1.0.0",
			"endpoints": fiber.Map{
				"health":     "GET /health",
				"create":     "POST /api/resources",
				"getAll":     "GET /api/resources?limit=10&offset=0",
				"getByID":    "GET /api/resources/:id",
				"getByCode":  "GET /api/resources/code/:code",
				"getByType":  "GET /api/resources/type/:type",
				"getByShift": "GET /api/resources/shift/:shiftId",
				"update":     "PUT /api/resources/:id",
				"delete":     "DELETE /api/resources/:id",
			},
		})
	})

	// Start server
	port := ":3000"
	log.Printf("🚀 Server starting on http://localhost%s", port)
	log.Printf("📚 API documentation: http://localhost%s/", port)
	log.Printf("❤️  Health check: http://localhost%s/health", port)
	log.Fatal(app.Listen(port))
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}

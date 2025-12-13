package server

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestLogRegisteredEndpoints(t *testing.T) {
	tests := []struct {
		name           string
		setupRoutes    func(*fiber.App)
		expectedLogs   []string
		unexpectedLogs []string
	}{
		{
			name: "logs registered endpoints",
			setupRoutes: func(app *fiber.App) {
				app.Get("/health", func(c *fiber.Ctx) error {
					return c.SendString("OK")
				})
				app.Post("/api/users", func(c *fiber.Ctx) error {
					return c.SendString("Created")
				})
			},
			expectedLogs: []string{
				"📋 Registered endpoints:",
				"[READY] GET /health",
				"[READY] POST /api/users",
			},
			unexpectedLogs: []string{
				"HEAD",
			},
		},
		{
			name:        "handles minimal endpoints",
			setupRoutes: func(app *fiber.App) {},
			expectedLogs: []string{
				"📋 Registered endpoints:",
				"[READY] GET /",
			},
		},
		{
			name: "filters out catch-all routes",
			setupRoutes: func(app *fiber.App) {
				app.Get("/test", func(c *fiber.Ctx) error {
					return c.SendString("OK")
				})
			},
			expectedLogs: []string{
				"[READY] GET /test",
			},
			unexpectedLogs: []string{
				"POST /",
				"PUT /",
				"DELETE /",
				"CONNECT /",
				"OPTIONS /",
				"TRACE /",
				"PATCH /",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(nil)

			// Create server with custom config
			config := DefaultConfig()
			config.DisableStartupMessage = true
			server := NewFiberServer(config)

			// Setup routes
			if tt.setupRoutes != nil {
				tt.setupRoutes(server.GetApp())
			}

			// Call the method under test
			server.LogRegisteredEndpoints()

			// Get logged output
			output := buf.String()

			// Check expected logs
			for _, expected := range tt.expectedLogs {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected log to contain %q, but it didn't. Output:\n%s", expected, output)
				}
			}

			// Check unexpected logs
			for _, unexpected := range tt.unexpectedLogs {
				if strings.Contains(output, unexpected) {
					t.Errorf("Expected log NOT to contain %q, but it did. Output:\n%s", unexpected, output)
				}
			}
		})
	}
}

func TestLogRegisteredEndpoints_MultipleMethodsSamePath(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	// Create server with custom config
	config := DefaultConfig()
	config.DisableStartupMessage = true
	server := NewFiberServer(config)

	// Setup routes with multiple methods on same path
	app := server.GetApp()
	app.Get("/api/resource", func(c *fiber.Ctx) error {
		return c.SendString("GET")
	})
	app.Post("/api/resource", func(c *fiber.Ctx) error {
		return c.SendString("POST")
	})
	app.Put("/api/resource", func(c *fiber.Ctx) error {
		return c.SendString("PUT")
	})

	// Call the method under test
	server.LogRegisteredEndpoints()

	// Get logged output
	output := buf.String()

	// Check that all methods are logged
	expectedMethods := []string{"GET", "POST", "PUT"}
	for _, method := range expectedMethods {
		expectedLog := "[READY] " + method + " /api/resource"
		if !strings.Contains(output, expectedLog) {
			t.Errorf("Expected log to contain %q, but it didn't. Output:\n%s", expectedLog, output)
		}
	}
}

func TestFiberServer_GetApp(t *testing.T) {
	config := DefaultConfig()
	server := NewFiberServer(config)

	app := server.GetApp()
	if app == nil {
		t.Error("Expected GetApp to return non-nil app")
	}
}

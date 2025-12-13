package main

import (
	"log"
	"os"
	"strconv"
	"time"

	userusecase "github.com/davidgaspardev/usermes-backend/internal/modules/iam/application/usecase"
	userhttp "github.com/davidgaspardev/usermes-backend/internal/modules/iam/infrastructure/adapter/input/http"
	userpersistence "github.com/davidgaspardev/usermes-backend/internal/modules/iam/infrastructure/adapter/output/persistence"
	resourceusecase "github.com/davidgaspardev/usermes-backend/internal/modules/production/application/usecase"
	resourcehttp "github.com/davidgaspardev/usermes-backend/internal/modules/production/infrastructure/adapter/input/http"
	resourcepersistence "github.com/davidgaspardev/usermes-backend/internal/modules/production/infrastructure/adapter/output/persistence"
	"github.com/davidgaspardev/usermes-backend/internal/shared/infrastructure/http/server"
	"github.com/davidgaspardev/usermes-backend/internal/shared/infrastructure/security"
)

const (
	// DefaultServerPort is the default port the server will listen on
	DefaultServerPort = 3001
)

func main() {
	// Load configuration (in a real app, use environment variables or config files)
	config := loadConfig()

	// Initialize infrastructure dependencies
	tokenGenerator := security.NewJWTTokenGenerator(config.JWTSecret, config.AppName)

	// Initialize User module
	userRepository := userpersistence.NewMemoryUserRepository()
	userService := userusecase.NewUserService(
		userRepository,
		tokenGenerator,
		config.TokenDuration,
	)

	// Initialize Resource module
	resourceRepository := resourcepersistence.NewMemoryResourceRepository()
	resourceService := resourceusecase.NewResourceService(resourceRepository)

	// Initialize HTTP server
	serverConfig := server.DefaultConfig()
	serverConfig.Port = config.ServerPort
	serverConfig.AppName = config.AppName

	httpServer := server.NewFiberServer(serverConfig)

	// Add health check endpoint
	httpServer.AddHealthCheck()

	// Register module routes
	// User module routes
	userRoutes := userhttp.NewUserRoutes(userService, tokenGenerator)
	httpServer.RegisterRoutes(userRoutes.SetupRoutes)

	// Resource module routes
	resourceRoutes := resourcehttp.NewResourceRoutes(resourceService)
	httpServer.RegisterRoutes(resourceRoutes.SetupRoutes)

	// Log all registered endpoints
	httpServer.LogRegisteredEndpoints()

	// Start server
	log.Printf("Starting %s on port %d", config.AppName, config.ServerPort)
	if err := httpServer.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Config holds the application configuration
type Config struct {
	AppName       string
	JWTSecret     string
	TokenDuration time.Duration
	ServerPort    int
}

// parseValidPort attempts to parse a port string and validates it's in the valid range (1-65535)
// Returns the parsed port and true if valid, otherwise returns 0 and false
func parseValidPort(portStr string) (int, bool) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, false
	}
	if port < 1 || port > 65535 {
		return 0, false
	}
	return port, true
}

// loadConfig loads the application configuration
// Reads from environment variables with fallback to default values
func loadConfig() Config {
	// Read port from environment variable, default to DefaultServerPort
	port := DefaultServerPort
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, valid := parseValidPort(portStr); valid {
			port = p
		} else {
			log.Printf("Warning: Invalid PORT value '%s' (must be 1-65535), using default port %d", portStr, DefaultServerPort)
		}
	}

	return Config{
		AppName:       "UserMes API",
		ServerPort:    port,
		JWTSecret:     "your-secret-key-change-this-in-production",
		TokenDuration: 24 * time.Hour, // 24 hours
	}
}

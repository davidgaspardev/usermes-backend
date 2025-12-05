package main

import (
	"log"
	"time"

	resourceusecase "github.com/davidgaspardev/usermes-backend/internal/modules/resource/application/usecase"
	resourcehttp "github.com/davidgaspardev/usermes-backend/internal/modules/resource/infrastructure/adapter/input/http"
	resourcepersistence "github.com/davidgaspardev/usermes-backend/internal/modules/resource/infrastructure/adapter/output/persistence"
	userusecase "github.com/davidgaspardev/usermes-backend/internal/modules/user/application/usecase"
	userhttp "github.com/davidgaspardev/usermes-backend/internal/modules/user/infrastructure/adapter/input/http"
	userpersistence "github.com/davidgaspardev/usermes-backend/internal/modules/user/infrastructure/adapter/output/persistence"
	"github.com/davidgaspardev/usermes-backend/internal/shared/infrastructure/http/server"
	"github.com/davidgaspardev/usermes-backend/internal/shared/infrastructure/security"
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

// loadConfig loads the application configuration
// In a production app, use environment variables or a config file
func loadConfig() Config {
	return Config{
		AppName:       "UserMes API",
		ServerPort:    3001,
		JWTSecret:     "your-secret-key-change-this-in-production",
		TokenDuration: 24 * time.Hour, // 24 hours
	}
}

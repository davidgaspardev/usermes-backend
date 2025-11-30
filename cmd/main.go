package main

import (
	"log"
	"time"

	"github.com/davidgaspardev/usermes-backend/internal/modules/user/application/usecase"
	userhttp "github.com/davidgaspardev/usermes-backend/internal/modules/user/infrastructure/adapter/input/http"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/infrastructure/adapter/output/persistence"
	"github.com/davidgaspardev/usermes-backend/internal/shared/infrastructure/http/server"
	"github.com/davidgaspardev/usermes-backend/internal/shared/infrastructure/security"
)

func main() {
	// Load configuration (in a real app, use environment variables or config files)
	config := loadConfig()

	// Initialize infrastructure dependencies
	tokenGenerator := security.NewJWTTokenGenerator(config.JWTSecret, config.AppName)
	userRepository := persistence.NewMemoryUserRepository()

	// Initialize application services (use cases)
	userService := usecase.NewUserService(
		userRepository,
		tokenGenerator,
		config.TokenDuration,
	)

	// Initialize HTTP server
	serverConfig := server.DefaultConfig()
	serverConfig.Port = config.ServerPort
	serverConfig.AppName = config.AppName

	httpServer := server.NewFiberServer(serverConfig)

	// Add health check endpoint
	httpServer.AddHealthCheck()

	// Register module routes
	userRoutes := userhttp.NewUserRoutes(userService, tokenGenerator)
	httpServer.RegisterRoutes(userRoutes.SetupRoutes)

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
		ServerPort:    3000,
		JWTSecret:     "your-secret-key-change-this-in-production",
		TokenDuration: 24 * time.Hour, // 24 hours
	}
}

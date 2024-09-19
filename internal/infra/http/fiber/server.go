package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	engine *fiber.App
	port   uint16
}

func NewServer(port uint16) Server {
	return Server{
		engine: fiber.New(),
		port:   port,
	}
}

func (srv *Server) Start() {
	// Register routes
	srv.registerRoutes()

	// Listen on port 3000
	srv.engine.Listen(fmt.Sprintf(":%d", srv.port))
}

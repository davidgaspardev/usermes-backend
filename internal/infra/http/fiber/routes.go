package server

import "github.com/gofiber/fiber/v2"

func (srv *Server) registerRoutes() {
	// Register your routes here
	srv.engine.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
}

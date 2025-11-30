package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/application/port/input"
)

// ResourceRoutes configures all routes for the resource module
type ResourceRoutes struct {
	handler *ResourceHandler
}

// NewResourceRoutes creates a new instance of ResourceRoutes
func NewResourceRoutes(resourceService input.ResourceService) *ResourceRoutes {
	return &ResourceRoutes{
		handler: NewResourceHandler(resourceService),
	}
}

// SetupRoutes registers all resource routes with the Fiber app
func (r *ResourceRoutes) SetupRoutes(app *fiber.App) {
	// Create a route group for resource endpoints
	resources := app.Group("/api/resources")

	// Create a new resource
	resources.Post("/", r.handler.Create)

	// Get all resources with pagination
	resources.Get("/", r.handler.GetAll)

	// Get resource by ID
	resources.Get("/:id", r.handler.GetByID)

	// Get resource by code
	resources.Get("/code/:code", r.handler.GetByCode)

	// Get resources by type
	resources.Get("/type/:type", r.handler.GetByType)

	// Get resources by shift ID
	resources.Get("/shift/:shiftId", r.handler.GetByShiftID)

	// Update a resource
	resources.Put("/:id", r.handler.Update)

	// Delete a resource
	resources.Delete("/:id", r.handler.Delete)
}

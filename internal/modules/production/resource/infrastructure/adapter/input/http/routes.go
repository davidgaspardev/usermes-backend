package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/resource/application/port/input"
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
// Routes are scoped to a specific plant: /v1/plants/:plant_code/production/resources
func (r *ResourceRoutes) SetupRoutes(app *fiber.App) {
	// Create a route group for plant-scoped production resources
	// Pattern: /v1/plants/:plant_code/production/resources
	resources := app.Group("/v1/plants/:plant_code/production/resources")

	// Create a new resource in the specified plant
	resources.Post("/", r.handler.Create)

	// Get all resources for the specified plant with pagination
	resources.Get("/", r.handler.GetAll)

	// Get resource by ID
	resources.Get("/:id", r.handler.GetByID)

	// Get resource by code within the specified plant
	resources.Get("/code/:code", r.handler.GetByCode)

	// Get resources by type within the specified plant
	resources.Get("/type/:type", r.handler.GetByType)

	// Get resources by shift ID within the specified plant
	resources.Get("/shift/:shiftId", r.handler.GetByShiftID)

	// Update a resource
	resources.Put("/:id", r.handler.Update)

	// Delete a resource
	resources.Delete("/:id", r.handler.Delete)
}

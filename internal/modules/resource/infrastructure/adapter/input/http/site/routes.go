package site

import (
	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/application/port/input"
	"github.com/gofiber/fiber/v2"
)

type SiteRoutes struct {
	handler *SiteHandler
}

func NewSiteRoutes(siteService input.SiteService) *SiteRoutes {
	return &SiteRoutes{
		handler: NewSiteHandler(siteService),
	}
}

func (r *SiteRoutes) SetupRoutes(app *fiber.App) {
	sites := app.Group("/api/sites")

	sites.Post("/", r.handler.Create)
}

package laneRoutes

import (
	handlers "bug-point-be/handlers/lanes"

	"github.com/gofiber/fiber/v2"
)

func SetupLaneRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	bugs := api.Group("/lanes")
	bugs.Get("/", handlers.GetLanes)
	bugs.Post("/", handlers.CreateLane)
	bugs.Delete("/:id", handlers.DeleteLane)
	bugs.Patch("/:id", handlers.Updatelane)
}

package bugRoutes

import (
	handlers "bug-point-be/handlers/bugs"

	"github.com/gofiber/fiber/v2"
)

func SetupBugRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	bugs := api.Group("/bugs")
	bugs.Get("/", handlers.GetBugs)
	bugs.Get("/:id", handlers.GetSingleBug)
	bugs.Post("/", handlers.CreateBug)
	bugs.Delete("/:id", handlers.DeleteBug)
	bugs.Put("/:id", handlers.UpdateBug)
}

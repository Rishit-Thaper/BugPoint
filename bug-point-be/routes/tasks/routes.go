package taskRoutes

import (
	handlers "bug-point-be/handlers/tasks"

	"github.com/gofiber/fiber/v2"
)

func SetupTaskRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	bugs := api.Group("/tasks")
	bugs.Get("/:laneId/", handlers.GetTasks)
	bugs.Get("/:id", handlers.GetSingleTask)
	bugs.Post("/:laneId/", handlers.CreateTask)
	bugs.Delete("/:id", handlers.DeleteTask)
	bugs.Patch("/:id", handlers.UpdateTask)
}

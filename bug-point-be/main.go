package main

import (
	config "bug-point-be/configs"
	db "bug-point-be/db"
	taskRoutes "bug-point-be/routes/tasks"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	client, ctx, cancel, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	app := fiber.New()
	taskRoutes.SetupTaskRoutes(app)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello world to the API"})
	})
	log.Fatal(app.Listen(":4000"))
}

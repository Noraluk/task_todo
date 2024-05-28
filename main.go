package main

import (
	"encoding/json"
	"errors"
	"todo/api/routes"
	"todo/pkg/config"
	"todo/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	err = database.Init()
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			return c.Status(code).JSON(e)
		},
	})
	app.Use(cors.New())

	routes.NewRoutes(app)
	app.Listen(":8080")
}

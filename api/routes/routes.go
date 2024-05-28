package routes

import (
	"todo/api/models/response"

	"github.com/gofiber/fiber/v2"
)

func NewRoutes(app *fiber.App) {
	handler := NewHandler()

	apiGroup := app.Group("/api")
	apiGroup.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(response.Response{Status: fiber.StatusOK})
	})

	taskGroup := apiGroup.Group("/tasks")
	taskGroup.Post("", handler.task.CreateTask)
	taskGroup.Get("", handler.task.GetTasks)
	taskGroup.Put("/:id", handler.task.UpdateTask)
}

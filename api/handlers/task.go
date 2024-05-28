package handlers

import (
	"todo/api/models/request"
	"todo/api/models/response"
	"todo/api/services"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler interface {
	CreateTask(c *fiber.Ctx) error
	GetTasks(c *fiber.Ctx) error
	UpdateTask(c *fiber.Ctx) error
}

type taskHandler struct {
	taskService services.TaskService
}

func NewTaskHandler(taskService services.TaskService) TaskHandler {
	return &taskHandler{
		taskService: taskService,
	}
}

func (h taskHandler) CreateTask(c *fiber.Ctx) error {
	var req request.CreatedTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err := req.Validate()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = h.taskService.CreateTask(req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(response.Response{
		Status: fiber.StatusOK,
	})
}

func (h taskHandler) GetTasks(c *fiber.Ctx) error {
	var query request.TaskListQuery
	if err := c.QueryParser(&query); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err := query.Validate()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tasks, err := h.taskService.GetTasks(query)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(response.Response{Status: fiber.StatusOK, Data: tasks})
}

func (h taskHandler) UpdateTask(c *fiber.Ctx) error {
	var req request.UpdatedTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err := req.Validate()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = h.taskService.UpdateTask(id, req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(response.Response{
		Status: fiber.StatusOK,
	})
}

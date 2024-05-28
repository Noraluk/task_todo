package routes

import (
	"todo/api/handlers"
	"todo/api/services"
	"todo/pkg/base"
	"todo/pkg/database"
)

type handler struct {
	task handlers.TaskHandler
}

func NewHandler() handler {
	repository := base.NewBaseRepository[any](database.GetDatabase())

	// services
	taskService := services.NewTaskService(repository)

	return handler{
		task: handlers.NewTaskHandler(taskService),
	}
}

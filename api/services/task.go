package services

import (
	"fmt"
	"time"
	"todo/api/entities"
	"todo/api/models/request"
	"todo/pkg/base"
	"todo/pkg/logger"
)

type TaskService interface {
	CreateTask(req request.CreatedTaskRequest) error
	GetTasks(query request.TaskListQuery) ([]entities.Task, error)
	UpdateTask(id int, req request.UpdatedTaskRequest) error
}

type taskService struct {
	repository base.BaseRepository[any]
	log        logger.Logger
}

func NewTaskService(repository base.BaseRepository[any]) TaskService {
	return &taskService{
		repository: repository,
		log:        logger.WithPrefix("service/task"),
	}
}

func (s taskService) CreateTask(req request.CreatedTaskRequest) error {
	tn := time.Now()
	task := entities.Task{
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   tn,
		UpdatedAt:   tn,
		Image:       req.Image,
		Status:      req.Status,
	}

	err := s.repository.Create(&task).Error()
	if err != nil {
		return err
	}
	return nil
}

func (s taskService) GetTasks(query request.TaskListQuery) ([]entities.Task, error) {
	var tasks []entities.Task
	db := s.repository

	if len(query.Title) > 0 {
		db = db.Where("title LIKE ?", fmt.Sprintf("%s%%", query.Title))
	}
	if len(query.Description) > 0 {
		db = db.Where("description LIKE ?", fmt.Sprintf("%s%%", query.Description))
	}
	if len(query.SortOrder) > 0 && len(query.SortBy) > 0 {
		db = db.Order(fmt.Sprintf("%s %s", query.SortBy, query.SortOrder))
	}

	err := db.Find(&tasks).Error()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s taskService) UpdateTask(id int, req request.UpdatedTaskRequest) error {
	updated := make(map[string]interface{})
	if len(req.Title) > 0 {
		updated["title"] = req.Title
	}
	if len(req.Description) > 0 {
		updated["description"] = req.Description
	}
	if len(req.Image) > 0 {
		updated["image"] = req.Image
	}
	if len(req.Status) > 0 {
		updated["status"] = req.Status
	}

	updated["updated_at"] = time.Now()
	err := s.repository.Model(&entities.Task{}).Where("id = ?", id).Updates(updated).Error()
	if err != nil {
		return err
	}

	return nil
}

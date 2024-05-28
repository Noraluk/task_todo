package request

import (
	"errors"
	"todo/api/enum"
)

type CreatedTaskRequest struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	Status      enum.TaskStatus `json:"status"`
}

func (r CreatedTaskRequest) Validate() error {
	if len(r.Title) > 100 {
		return errors.New("title is exceeded more than 100")
	}

	if !r.Status.IsValid() {
		return errors.New("status is invalid")
	}

	return nil
}

type TaskListQuery struct {
	Title       string              `query:"title"`
	Description string              `query:"description"`
	SortBy      enum.TaskListSortBy `query:"sort_by"`
	SortOrder   enum.SortOrder      `query:"sort_order"`
}

func (r TaskListQuery) Validate() error {
	if len(r.SortBy) > 0 && !r.SortBy.IsValid() {
		return errors.New("sort by is invalid")
	}

	if len(r.SortOrder) > 0 && !r.SortOrder.IsValid() {
		return errors.New("sort order is invalid")
	}

	return nil
}

type UpdatedTaskRequest struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	Status      enum.TaskStatus `json:"status"`
}

func (r UpdatedTaskRequest) Validate() error {
	if len(r.Title) > 100 {
		return errors.New("title is exceeded more than 100")
	}

	if len(r.Status) > 0 && !r.Status.IsValid() {
		return errors.New("status is invalid")
	}

	return nil
}

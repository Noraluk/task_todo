package entities

import (
	"time"
	"todo/api/enum"
)

type Task struct {
	ID          int             `gorm:"primaryKey" json:"id"`
	Title       string          `gorm:"size:100" json:"title"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Image       string          `json:"image"`
	Status      enum.TaskStatus `json:"status"`
}

package enum

type TaskStatus string

const (
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
)

func (e TaskStatus) IsValid() bool {
	switch e {
	case TaskStatusInProgress, TaskStatusCompleted:
		return true
	}
	return false
}

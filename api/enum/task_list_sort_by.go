package enum

type TaskListSortBy string

const (
	TaskListSortByTitle     TaskListSortBy = "title"
	TaskListSortByCreatedAt                = "created_at"
	TaskListSortByUpdatedAt                = "updated_at"
	TaskListSortByStatus    TaskListSortBy = "status"
)

func (e TaskListSortBy) IsValid() bool {
	switch e {
	case TaskListSortByTitle, TaskListSortByCreatedAt, TaskListSortByUpdatedAt, TaskListSortByStatus:
		return true
	}
	return false
}

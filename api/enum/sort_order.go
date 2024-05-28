package enum

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

func (e SortOrder) IsValid() bool {
	switch e {
	case SortOrderAsc, SortOrderDesc:
		return true
	}
	return false
}

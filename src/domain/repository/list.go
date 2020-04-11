package repository

// Sort :
type Sort struct {
	By    SortBy
	Order SortOrder
}

// SortBy :
type SortBy string

var (
	// SortByTaskRemains :
	SortByTaskRemains = SortBy("TaskRemains")
	// SortByTaskDoneAt :
	SortByTaskDoneAt = SortBy("TaskDoneAt")
)

// SortOrder : asc or desc
type SortOrder string

var (
	// SortOrderAsc : ascending order
	SortOrderAsc = SortOrder("ASC")
	// SortOrderDesc : decending order
	SortOrderDesc = SortOrder("DESC")
)

// ListOption :
type ListOption struct {
	Sort   Sort
	Limit  int
	Offset int
}

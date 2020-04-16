package sqliteimpl

import (
	"fmt"

	"github.com/notomo/counteria.nvim/src/domain/repository"
)

func convertListOption(option repository.ListOption) string {
	// TODO : limit, offset
	return convertSort(option.Sort)
}

func convertSort(sort repository.Sort) string {
	var by string
	switch sort.By {
	case repository.SortByTaskDoneAt:
		by = "done.at"
	case repository.SortByTaskRemains:
		return ""
	default:
		panic("invalid sort by: " + sort.By)
	}
	return fmt.Sprintf("ORDER BY %s %s", by, sort.Order)
}

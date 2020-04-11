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
	default:
		return ""
	}
	return fmt.Sprintf("ORDER BY %s %s", by, sort.Order)
}

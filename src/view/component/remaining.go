package component

import (
	"fmt"

	"github.com/notomo/counteria.nvim/src/domain/model"
)

// RemainingTime : view
type RemainingTime struct {
	model.RemainingTime
}

func (time RemainingTime) String() string {
	var sign string
	if !time.Exists() {
		sign = "- "
	}

	var days string
	if time.Days != 0 {
		days = fmt.Sprintf("%d days ", time.Days)
	}

	var hours string
	if days != "" || time.Hours != 0 {
		hours = fmt.Sprintf("%d hours ", time.Hours)
	}

	var minutes string
	if hours != "" || time.Minutes != 0 {
		minutes = fmt.Sprintf("%d minutes", time.Minutes)
	}

	return fmt.Sprintf("%s%s%s%s", sign, days, hours, minutes)
}

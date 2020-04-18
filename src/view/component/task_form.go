package component

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/pkg/errors"
)

// NewTaskForm :
func NewTaskForm(task *model.Task) *TaskFormView {
	rule := task.Rule()

	periods := []PeriodView{}
	for _, p := range rule.Periods() {
		periods = append(periods, PeriodView{
			PeriodNumber: p.Number(),
			PeriodUnit:   p.Unit(),
		})
	}

	return &TaskFormView{
		TaskName:    task.Name(),
		TaskStartAt: task.StartAt(),
		TaskRuleView: TaskRuleView{
			RuleType:      rule.Type(),
			RuleWeekdays:  rule.Weekdays(),
			RuleDays:      rule.Days(),
			RuleMonthDays: rule.MonthDays(),
			RuleDateTimes: rule.DateTimes(),
			RuleDates:     rule.Dates(),
			RulePeriods:   periods,
		},
	}
}

var _ model.PeriodData = PeriodView{}

// PeriodView :
type PeriodView struct {
	PeriodNumber int              `json:"number"`
	PeriodUnit   model.PeriodUnit `json:"unit"`
}

// Number :
func (view PeriodView) Number() int {
	return view.PeriodNumber
}

// Unit :
func (view PeriodView) Unit() model.PeriodUnit {
	return view.PeriodUnit
}

// TaskFormView :
type TaskFormView struct {
	TaskID      int       `json:"-"`
	TaskName    string    `json:"name"`
	TaskStartAt time.Time `json:"startAt"`
	TaskRuleView
}

// Lines :
func (view *TaskFormView) Lines() ([][]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(&view); err != nil {
		return nil, errors.WithStack(err)
	}
	lines := bytes.Split(b.Bytes(), []byte("\n"))
	return lines[:len(lines)-1], nil
}

var _ model.TaskRuleData = &TaskRuleView{}

// TaskRuleView :
type TaskRuleView struct {
	RuleType      model.TaskRuleType `json:"type"`
	RuleWeekdays  model.Weekdays     `json:"weekdays"`
	RuleDays      model.Days         `json:"days"`
	RuleMonthDays model.MonthDays    `json:"monthDays"`
	RuleDateTimes model.DateTimes    `json:"dateTimes"`
	RuleDates     model.Dates        `json:"dates"`
	RulePeriods   []PeriodView       `json:"periods"`
}

// Type :
func (view *TaskRuleView) Type() model.TaskRuleType {
	return view.RuleType
}

// Weekdays :
func (view *TaskRuleView) Weekdays() model.Weekdays {
	return view.RuleWeekdays
}

// Days :
func (view *TaskRuleView) Days() model.Days {
	return view.RuleDays
}

// MonthDays :
func (view *TaskRuleView) MonthDays() model.MonthDays {
	return view.RuleMonthDays
}

// Dates :
func (view *TaskRuleView) Dates() model.Dates {
	return view.RuleDates
}

// DateTimes :
func (view *TaskRuleView) DateTimes() model.DateTimes {
	return view.RuleDateTimes
}

// Periods :
func (view *TaskRuleView) Periods() model.Periods {
	periods := model.Periods{}
	for _, p := range view.RulePeriods {
		periods = append(periods, model.Period{
			PeriodData: p,
		})
	}
	return periods
}

var _ model.TaskData = &TaskFormView{}

// ID :
func (view *TaskFormView) ID() int {
	return view.TaskID
}

// Name :
func (view *TaskFormView) Name() string {
	return view.TaskName
}

// Rule :
func (view *TaskFormView) Rule() *model.TaskRule {
	return &model.TaskRule{
		TaskRuleData: &view.TaskRuleView,
	}
}

// StartAt :
func (view *TaskFormView) StartAt() time.Time {
	return view.TaskStartAt
}

// LastDone :
func (view *TaskFormView) LastDone() *model.DoneTask {
	return nil
}

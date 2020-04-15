package view

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/pkg/errors"
)

// PeriodView :
type PeriodView struct {
	PeriodNumber int              `json:"number"`
	PeriodUnit   model.PeriodUnit `json:"unit"`
}

// Number :
func (period PeriodView) Number() int {
	return period.PeriodNumber
}

// Unit :
func (period PeriodView) Unit() model.PeriodUnit {
	return period.PeriodUnit
}

// TaskFormView :
type TaskFormView struct {
	TaskID      int       `json:"-"`
	TaskName    string    `json:"name"`
	TaskStartAt time.Time `json:"startAt"`
	TaskRuleView
}

// TaskRuleView :
type TaskRuleView struct {
	RuleWeekdays  model.Weekdays  `json:"weekdays"`
	RuleDays      model.Days      `json:"days"`
	RuleMonthDays model.MonthDays `json:"monthDays"`
	RuleTimes     model.Times     `json:"times"`
	RuleDates     model.Dates     `json:"dates"`
	RulePeriods   []PeriodView    `json:"periods"`
}

// Type :
func (rule *TaskRuleView) Type() model.TaskRuleType {
	switch {
	case len(rule.RuleWeekdays) != 0:
		return model.TaskRuleTypeInWeekdays
	case len(rule.RuleDays) != 0:
		return model.TaskRuleTypeInDaysEveryMonth
	case len(rule.RuleMonthDays) != 0:
		return model.TaskRuleTypeInDaysEveryMonth
	case len(rule.RuleTimes) != 0:
		return model.TaskRuleTypeByTimes
	case len(rule.RuleDates) != 0:
		return model.TaskRuleTypeInDates
	case len(rule.RulePeriods) != 0:
		return model.TaskRuleTypePeriodic
	}
	panic("invalid rule")
}

// Weekdays :
func (rule *TaskRuleView) Weekdays() model.Weekdays {
	return rule.RuleWeekdays
}

// Days :
func (rule *TaskRuleView) Days() model.Days {
	return rule.RuleDays
}

// MonthDays :
func (rule *TaskRuleView) MonthDays() model.MonthDays {
	return rule.RuleMonthDays
}

// Dates :
func (rule *TaskRuleView) Dates() model.Dates {
	return rule.RuleDates
}

// DateTimes :
func (rule *TaskRuleView) DateTimes() model.Times {
	return rule.RuleTimes
}

// Periods :
func (rule *TaskRuleView) Periods() model.Periods {
	periods := model.Periods{}
	for _, p := range rule.RulePeriods {
		periods = append(periods, model.Period{
			PeriodData: p,
		})
	}
	return periods
}

var _ model.TaskData = &TaskFormView{}

// ID :
func (task *TaskFormView) ID() int {
	return task.TaskID
}

// Name :
func (task *TaskFormView) Name() string {
	return task.TaskName
}

// Rule :
func (task *TaskFormView) Rule() *model.TaskRule {
	return &model.TaskRule{
		TaskRuleData: &task.TaskRuleView,
	}
}

// StartAt :
func (task *TaskFormView) StartAt() time.Time {
	return task.TaskStartAt
}

// LastDone :
func (task *TaskFormView) LastDone() *model.DoneTask {
	return nil
}

// TaskFromForm :
func (renderer *BufferRenderer) TaskFromForm(id int) (*model.Task, error) {
	reader, err := renderer.Buffer.Reader()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var view TaskFormView
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&view); err != nil {
		return nil, errors.WithStack(err)
	}
	view.TaskID = id

	return &model.Task{
		TaskData: &view,
	}, nil
}

// OneNewTask :
func (renderer *BufferRenderer) OneNewTask(task *model.Task) error {
	rule := task.Rule()

	periods := []PeriodView{}
	for _, p := range rule.Periods() {
		periods = append(periods, PeriodView{
			PeriodNumber: p.Number(),
			PeriodUnit:   p.Unit(),
		})
	}

	view := &TaskFormView{
		TaskName:    task.Name(),
		TaskStartAt: task.StartAt(),
		TaskRuleView: TaskRuleView{
			RuleWeekdays:  rule.Weekdays(),
			RuleDays:      rule.Days(),
			RuleMonthDays: rule.MonthDays(),
			RuleTimes:     rule.DateTimes(),
			RuleDates:     rule.Dates(),
			RulePeriods:   periods,
		},
	}

	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(view); err != nil {
		return errors.WithStack(err)
	}
	lines := bytes.Split(b.Bytes(), []byte("\n"))

	if err := renderer.Buffer.SetLines(
		lines[:len(lines)-1],
		renderer.Buffer.WithBufferType("acwrite"),
		renderer.Buffer.WithModifiable(true),
		renderer.Buffer.WithOpen(),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// OneTask :
func (renderer *BufferRenderer) OneTask(task *model.Task) error {
	rule := task.Rule()

	periods := []PeriodView{}
	for _, p := range rule.Periods() {
		periods = append(periods, PeriodView{
			PeriodNumber: p.Number(),
			PeriodUnit:   p.Unit(),
		})
	}

	view := &TaskFormView{
		TaskName:    task.Name(),
		TaskStartAt: task.StartAt(),
		TaskRuleView: TaskRuleView{
			RuleWeekdays:  rule.Weekdays(),
			RuleDays:      rule.Days(),
			RuleMonthDays: rule.MonthDays(),
			RuleTimes:     rule.DateTimes(),
			RuleDates:     rule.Dates(),
			RulePeriods:   periods,
		},
	}

	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(&view); err != nil {
		return errors.WithStack(err)
	}
	lines := bytes.Split(b.Bytes(), []byte("\n"))

	if err := renderer.Buffer.SetLines(
		lines[:len(lines)-1],
		renderer.Buffer.WithBufferType("acwrite"),
		renderer.Buffer.WithModifiable(true),
		renderer.Buffer.WithOpen(),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

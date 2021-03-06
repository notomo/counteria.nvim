package view

import (
	"time"

	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/view/component"
	"github.com/notomo/counteria.nvim/src/vimlib"
	"github.com/pkg/errors"
)

func toLines(tasks []model.Task, now time.Time) ([][]byte, []vimlib.Highlight, error) {
	table, err := component.NewTable("", "Name", "Done", "Rule", "Remains")
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	highlights := []vimlib.Highlight{}
	for i, task := range tasks {
		at := "---------- --:--:--"
		doneAt := task.DoneAt()
		if doneAt != nil {
			at = doneAt.Format("2006-01-02 15:04:05")
		}
		deadline := task.Deadline(now)
		remainingTime := component.RemainingTime{RemainingTime: deadline.RemainingTime()}
		remaining := remainingTime.String()

		status := " "
		if !remainingTime.Exists() {
			status = "!"
			highlights = append(highlights, vimlib.Highlight{
				Group:    "Todo",
				Line:     i + 1,
				StartCol: 0,
				EndCol:   1,
			})
		}

		rule := task.Rule().String()
		if err := table.AddLine(status, task.Name(), at, rule, remaining); err != nil {
			return nil, nil, errors.WithStack(err)
		}
	}

	lines, highs, err := table.Lines(
		table.WithColumnHighlightGroup("TabLineSel"),
	)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return lines, append(highs, highlights...), nil
}

// TaskList :
func (renderer *BufferRenderer) TaskList(tasks []model.Task, now time.Time) error {
	lines, highlights, err := toLines(tasks, now)
	if err != nil {
		return errors.WithStack(err)
	}

	markIDs := make([]int, len(tasks))
	if err := renderer.Buffer.SetLines(
		lines,
		renderer.Buffer.WithBufferType("nofile"),
		renderer.Buffer.WithFileType("counteria-tasks"),
		renderer.Buffer.WithModifiable(false),
		renderer.Buffer.WithExtmarks(markIDs, 1),
		renderer.Buffer.WithHighlights(highlights),
	); err != nil {
		return errors.WithStack(err)
	}

	states := vimlib.LineStates{}
	for i, task := range tasks {
		path, err := route.TasksOnePath(task.ID())
		if err != nil {
			return errors.WithStack(err)
		}
		states.Add(markIDs[i], path)
	}
	if err := renderer.Buffer.SaveLineState(states); err != nil {
		return errors.WithStack(err)
	}

	if err := renderer.Buffer.Open(
		renderer.Buffer.WithWindowOption("list", false),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

package view

import (
	"bytes"
	"fmt"
	"time"

	"github.com/WeiZhang555/tabwriter"
	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/vimlib"
	"github.com/pkg/errors"
)

func toLines(tasks []model.Task, now time.Time) ([][]byte, []vimlib.Highlight, error) {
	var b bytes.Buffer
	minwidth, tabwidth := 1, 1
	padding := 2
	noflag := uint(0)
	w := tabwriter.NewWriter(&b, minwidth, tabwidth, padding, ' ', noflag)
	w.Write([]byte("\tName\tDone\tRule\tRemains\n"))

	highlights := []vimlib.Highlight{
		{
			Group:    "TabLineSel",
			Line:     0,
			StartCol: 0,
			EndCol:   -1,
		},
	}
	for i, task := range tasks {
		period := task.Period()

		at := "---------- --:--:--"
		doneAt := task.DoneAt()
		if doneAt != nil {
			at = doneAt.Format("2006-01-02 15:04:05")
		}
		remainingTime := task.RemainingTime(now)
		remaining := fmt.Sprintf("%d days %d hours %d minutes", remainingTime.Days, remainingTime.Hours, remainingTime.Minutes)

		status := " "
		if task.PastDeadline(now) {
			status = "!"
			highlights = append(highlights, vimlib.Highlight{
				Group:    "Todo",
				Line:     i + 1,
				StartCol: 0,
				EndCol:   1,
			})
		}

		line := fmt.Sprintf("%s\t%s\t%s\tonce per %d %s\t%s\n", status, task.Name(), at, period.Number(), period.Unit(), remaining)
		w.Write([]byte(line))
	}
	if err := w.Flush(); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	lines := bytes.Split(b.Bytes(), []byte("\n"))
	return lines[:len(lines)-1], highlights, nil
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

	if err := renderer.Buffer.Open(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

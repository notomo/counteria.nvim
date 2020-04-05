package view

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/WeiZhang555/tabwriter"
	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/vimlib"
	"github.com/pkg/errors"
)

func toLines(tasks []model.Task) ([][]byte, error) {
	var b bytes.Buffer
	minwidth, tabwidth := 1, 1
	padding := 4
	noflag := uint(0)
	w := tabwriter.NewWriter(&b, minwidth, tabwidth, padding, ' ', noflag)

	for _, task := range tasks {
		period := task.Period()

		at := "---------- --:--:--"
		doneAt := task.DoneAt()
		if doneAt != nil {
			at = doneAt.Format("2006-01-02 15:04:05")
		}
		limit := task.LimitAt().Format("2006-01-02 15:04:05")

		line := fmt.Sprintf("%s\t%s\tonce per %d %s\t%s\n", task.Name(), at, period.Number(), period.Unit(), limit)
		w.Write([]byte(line))
	}
	if err := w.Flush(); err != nil {
		return nil, errors.WithStack(err)
	}

	lines := bytes.Split(b.Bytes(), []byte("\n"))
	return lines[:len(lines)-1], nil
}

// TaskList :
func (renderer *BufferRenderer) TaskList(tasks []model.Task) error {
	lines, err := toLines(tasks)
	if err != nil {
		return errors.WithStack(err)
	}

	markIDs := make([]int, len(tasks))
	if err := renderer.Buffer.SetLines(
		lines,
		renderer.Buffer.WithBufferType("nofile"),
		renderer.Buffer.WithFileType("counteria-tasks"),
		renderer.Buffer.WithModifiable(false),
		renderer.Buffer.WithExtmarks(markIDs),
	); err != nil {
		return errors.WithStack(err)
	}

	states := vimlib.LineStates{}
	for i, task := range tasks {
		id := strconv.Itoa(markIDs[i])

		r := route.TasksOne
		params := route.Params{"taskId": strconv.Itoa(task.ID())}
		path, err := r.BuildPath(params)
		if err != nil {
			return errors.WithStack(err)
		}

		state := vimlib.LineState{Path: path}
		states[id] = state
	}
	if err := renderer.Buffer.SaveLineState(states); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

package view

import (
	"strconv"

	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/vimlib"
	"github.com/pkg/errors"
)

// TaskList :
func (renderer *BufferRenderer) TaskList(tasks []model.Task) error {
	lines := [][]byte{}
	for _, task := range tasks {
		lines = append(lines, []byte(task.Name()))
	}
	markIDs := make([]int, len(tasks))

	buffer := renderer.BufferClient
	if err := buffer.SetLines(
		lines,
		buffer.WithBufferType("nofile"),
		buffer.WithFileType("counteria-tasks"),
		buffer.WithModifiable(false),
		buffer.WithExtmarks(markIDs),
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
	if err := buffer.SaveLineState(states); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

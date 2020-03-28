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
		lines = append(lines, []byte(task.Name))
	}

	buf := renderer.Buffer
	batch := renderer.Vim.NewBatch()
	batch.ClearBufferNamespace(buf, renderer.NsID, 0, -1)
	batch.SetBufferOption(buf, "modifiable", true)
	batch.SetBufferLines(buf, 0, -1, false, lines)

	noneOpts := map[string]interface{}{}
	markIDs := make([]int, len(tasks))
	for i := range tasks {
		batch.SetBufferExtmark(buf, renderer.NsID, 0, i, 0, noneOpts, &markIDs[i])
	}

	batch.SetBufferOption(buf, "modifiable", false)
	batch.SetBufferOption(buf, "buftype", "nofile")
	batch.SetBufferOption(buf, "modified", false)
	batch.SetBufferOption(buf, "filetype", "counteria-tasks")
	if err := batch.Execute(); err != nil {
		return errors.WithStack(err)
	}

	bufState := map[string]vimlib.LineState{}
	for i, task := range tasks {
		id := strconv.Itoa(markIDs[i])

		r := route.TasksOne
		params := route.Params{"taskId": strconv.Itoa(task.ID)}
		path, err := r.BuildPath(params)
		if err != nil {
			return errors.WithStack(err)
		}

		state := vimlib.LineState{Path: path}
		bufState[id] = state
	}
	err := renderer.Vim.SetBufferVar(buf, "_counteria_state", bufState)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

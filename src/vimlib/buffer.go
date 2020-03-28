package vimlib

import (
	"strconv"
	"sync"

	"github.com/neovim/go-client/nvim"
	"github.com/pkg/errors"
)

// BufferClient :
type BufferClient struct {
	Vim *nvim.Nvim

	NsID         int
	getNamespace sync.Once
}

// LineState :
type LineState struct {
	Path string
}

// LineStates :
type LineStates map[string]LineState

const stateKeyName = "_counteria_state"

// SaveLineState :
func (client *BufferClient) SaveLineState(buf nvim.Buffer, states LineStates) error {
	if err := client.Vim.SetBufferVar(buf, stateKeyName, states); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// LoadLineState :
func (client *BufferClient) LoadLineState() (*LineState, error) {
	client.getNamespace.Do(func() {
		ns, err := client.Vim.CreateNamespace("counteria")
		if err != nil {
			panic(err)
		}
		client.NsID = ns
	})

	states := LineStates{}
	if err := client.Vim.BufferVar(0, stateKeyName, states); err != nil {
		return nil, errors.WithStack(err)
	}

	pos, err := client.Vim.WindowCursor(0)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	line := pos[0] - 1
	start := []int{line, 0}
	end := []int{line, -1}
	noneOpts := map[string]interface{}{}
	marks, err := client.Vim.BufferExtmarks(0, client.NsID, start, end, noneOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(marks) == 0 {
		return nil, nil
	}
	id := marks[0].ExtmarkID

	state, ok := states[strconv.Itoa(id)]
	if !ok {
		return nil, nil
	}

	return &state, nil
}

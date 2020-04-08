package vimlib

import (
	"github.com/neovim/go-client/nvim"
	"github.com/pkg/errors"
)

// SaveCursor :
func (client *BufferClient) SaveCursor() (*BufferCursor, error) {
	var path string
	var pos [2]int
	batch := client.Vim.NewBatch()
	batch.BufferName(client.Bufnr, &path)
	batch.WindowCursor(0, &pos)
	if err := batch.Execute(); err != nil {
		return nil, errors.WithStack(err)
	}

	return &BufferCursor{
		Vim:      client.Vim,
		Path:     path,
		Position: pos,
	}, nil
}

// BufferCursor :
type BufferCursor struct {
	Vim      *nvim.Nvim
	Path     string
	Position [2]int
}

// Restore :
func (cursor *BufferCursor) Restore() error {
	var path string
	var line int
	batch := cursor.Vim.NewBatch()
	batch.BufferName(0, &path)
	batch.BufferLineCount(0, &line)
	if err := batch.Execute(); err != nil {
		return errors.WithStack(err)
	}

	if cursor.Path != path {
		return nil
	}

	pos := cursor.Position
	if pos[0] > line {
		pos[0] = line
	}

	if err := cursor.Vim.SetWindowCursor(0, pos); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

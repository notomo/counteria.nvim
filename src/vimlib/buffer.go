package vimlib

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/neovim/go-client/nvim"
	"github.com/pkg/errors"
)

// BufferClientFactory :
type BufferClientFactory struct {
	Vim *nvim.Nvim

	NsID         int
	getNamespace sync.Once
}

// Current :
func (factory *BufferClientFactory) Current() *BufferClient {
	return factory.Get(0)
}

// GetOrCreate :
// NOTE: avoid executing BufReadCmd
func (factory *BufferClientFactory) GetOrCreate(path string) (*BufferClient, error) {
	var bufnr int
	pattern := fmt.Sprintf("^%s$", path)
	if err := factory.Vim.Call("bufnr", &bufnr, pattern); err != nil {
		return nil, errors.WithStack(err)
	}

	buf := nvim.Buffer(bufnr)
	exists := bufnr != -1
	if !exists {
		b, err := factory.Vim.CreateBuffer(false, true)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		buf = b
		if err := factory.Vim.SetBufferName(buf, path); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return factory.Get(buf), nil
}

// Get :
func (factory *BufferClientFactory) Get(bufnr nvim.Buffer) *BufferClient {
	factory.getNamespace.Do(func() {
		ns, err := factory.Vim.CreateNamespace("counteria")
		if err != nil {
			panic(err)
		}
		factory.NsID = ns
	})
	return &BufferClient{
		Vim:   factory.Vim,
		Bufnr: bufnr,
		NsID:  factory.NsID,
	}
}

// BufferClient :
type BufferClient struct {
	Vim   *nvim.Nvim
	Bufnr nvim.Buffer
	NsID  int
}

// LineState :
type LineState struct {
	Path string
}

// LineStates :
type LineStates map[string]LineState

// Add : to state
func (states LineStates) Add(markID int, path string) {
	id := strconv.Itoa(markID)
	states[id] = LineState{Path: path}
}

const stateKeyName = "_counteria_state"

// Open :
func (client *BufferClient) Open(opts ...func(*nvim.Batch)) error {
	batch := client.Vim.NewBatch()
	client.WithOpen()(batch)
	for _, opt := range opts {
		opt(batch)
	}

	if err := batch.Execute(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// WithWindowOption :
// NOTICE: this optin is valid only after window opened
func (client *BufferClient) WithWindowOption(name string, value interface{}) func(*nvim.Batch) {
	return func(batch *nvim.Batch) {
		batch.SetWindowOption(0, name, value)
	}
}

// SaveLineState :
func (client *BufferClient) SaveLineState(states LineStates) error {
	if err := client.Vim.SetBufferVar(client.Bufnr, stateKeyName, states); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// LoadLineState :
func (client *BufferClient) LoadLineState() (*LineState, error) {
	states := LineStates{}
	if err := client.Vim.BufferVar(client.Bufnr, stateKeyName, states); err != nil {
		return nil, ErrNoState
	}

	pos, err := client.Vim.WindowCursor(0)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	line := pos[0] - 1
	start := []int{line, 0}
	end := []int{line, -1}
	noneOpts := map[string]interface{}{}
	marks, err := client.Vim.BufferExtmarks(client.Bufnr, client.NsID, start, end, noneOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(marks) == 0 {
		return nil, ErrNoState
	}
	id := marks[0].ExtmarkID

	state, ok := states[strconv.Itoa(id)]
	if !ok {
		return nil, ErrNoState
	}

	return &state, nil
}

// Reader :
func (client *BufferClient) Reader() (io.Reader, error) {
	b, err := client.Vim.BufferLines(client.Bufnr, 0, -1, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return bytes.NewReader(bytes.Join(b, nil)), nil
}

// Save :
func (client *BufferClient) Save() error {
	batch := client.Vim.NewBatch()
	batch.SetBufferOption(client.Bufnr, "modified", false)
	if err := batch.Execute(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// SetLines :
func (client *BufferClient) SetLines(lines [][]byte, opts ...func(*nvim.Batch)) error {
	batch := client.Vim.NewBatch()
	batch.ClearBufferNamespace(client.Bufnr, client.NsID, 0, -1)
	batch.SetBufferOption(client.Bufnr, "modifiable", true)
	batch.SetBufferLines(client.Bufnr, 0, -1, false, lines)

	for _, opt := range opts {
		opt(batch)
	}

	batch.SetBufferOption(client.Bufnr, "modified", false)

	if err := batch.Execute(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// WithBufferType :
func (client *BufferClient) WithBufferType(typ string) func(*nvim.Batch) {
	return func(batch *nvim.Batch) {
		batch.SetBufferOption(client.Bufnr, "buftype", typ)
	}
}

// WithFileType :
func (client *BufferClient) WithFileType(typ string) func(*nvim.Batch) {
	return func(batch *nvim.Batch) {
		batch.SetBufferOption(client.Bufnr, "filetype", typ)
	}
}

// WithModifiable :
func (client *BufferClient) WithModifiable(modifiable bool) func(*nvim.Batch) {
	return func(batch *nvim.Batch) {
		batch.SetBufferOption(client.Bufnr, "modifiable", modifiable)
	}
}

// WithOpen :
func (client *BufferClient) WithOpen() func(*nvim.Batch) {
	cmd := fmt.Sprintf("buffer %d", client.Bufnr)
	return func(batch *nvim.Batch) {
		batch.Command(cmd)
	}
}

// WithExtmarks :
func (client *BufferClient) WithExtmarks(results []int, startLine int) func(*nvim.Batch) {
	noneOpts := map[string]interface{}{}
	return func(batch *nvim.Batch) {
		for i := range results {
			batch.SetBufferExtmark(client.Bufnr, client.NsID, 0, i+startLine, 0, noneOpts, &results[i])
		}
	}
}

// Highlight : color and position
type Highlight struct {
	Group    string
	Line     int
	StartCol int
	EndCol   int
}

// WithHighlights :
func (client *BufferClient) WithHighlights(highlights []Highlight) func(*nvim.Batch) {
	return func(batch *nvim.Batch) {
		for _, h := range highlights {
			var result int
			batch.AddBufferHighlight(client.Bufnr, client.NsID, h.Group, h.Line, h.StartCol, h.EndCol, &result)
		}
	}
}

// SyncRequest :
func (client *BufferClient) SyncRequest(method, path string) error {
	var unused interface{}
	if err := client.Vim.Call("counteria#request", unused, method, true, path, client.Bufnr); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Path :
func (client *BufferClient) Path() (string, error) {
	path, err := client.Vim.BufferName(client.Bufnr)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return path, nil
}

// BufferLinesEvent : nvim_buf_lines_event
type BufferLinesEvent struct {
	Bufnr       nvim.Buffer
	ChangedTick int
	FirstLine   int
	LastLine    int
	Lines       []string
	More        bool
}

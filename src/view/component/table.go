package component

import (
	"bytes"
	"strings"

	"github.com/WeiZhang555/tabwriter"
	"github.com/notomo/counteria.nvim/src/vimlib"
	"github.com/pkg/errors"
)

// Table :
type Table struct {
	Writer  *tabwriter.Writer
	Buffer  *bytes.Buffer
	Columns []string
}

// NewTable :
func NewTable(columns ...string) (*Table, error) {
	var b bytes.Buffer
	minwidth, tabwidth := 1, 1
	padding := 2
	noflag := uint(0)
	w := tabwriter.NewWriter(&b, minwidth, tabwidth, padding, ' ', noflag)
	if _, err := w.Write([]byte(strings.Join(columns, "\t|\t") + "\t|\n")); err != nil {
		return nil, err
	}
	return &Table{
		Writer:  w,
		Buffer:  &b,
		Columns: columns,
	}, nil
}

// AddLine :
func (table *Table) AddLine(cells ...string) error {
	if _, err := table.Writer.Write([]byte(strings.Join(cells, "\t|\t") + "\t|\n")); err != nil {
		return err
	}
	return nil
}

// LineOption : for table.Lines
type LineOption struct {
	// highlight group for all column labels
	ColumnHighlightGroup string
}

// Lines :
func (table *Table) Lines(opts ...func(*LineOption)) ([][]byte, []vimlib.Highlight, error) {
	option := &LineOption{}
	for _, opt := range opts {
		opt(option)
	}

	if err := table.Writer.Flush(); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	lines := bytes.Split(table.Buffer.Bytes(), []byte("\n"))
	highlights := table.columnHighlights(lines[0], option)

	return lines[:len(lines)-1], highlights, nil
}

func (table *Table) columnHighlights(line []byte, option *LineOption) []vimlib.Highlight {
	if option.ColumnHighlightGroup == "" {
		return nil
	}

	highlights := []vimlib.Highlight{}
	for _, column := range table.Columns {
		start := bytes.Index(line, []byte(column))
		if start == -1 {
			continue
		}
		highlights = append(highlights, vimlib.Highlight{
			Group:    option.ColumnHighlightGroup,
			Line:     0,
			StartCol: start,
			EndCol:   start + len(column),
		})
	}
	return highlights
}

// WithColumnHighlightGroup :
func (table *Table) WithColumnHighlightGroup(group string) func(*LineOption) {
	return func(op *LineOption) {
		op.ColumnHighlightGroup = group
	}
}

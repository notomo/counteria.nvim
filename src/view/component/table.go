package component

import (
	"bytes"
	"strings"

	"github.com/WeiZhang555/tabwriter"
	"github.com/pkg/errors"
)

// Table :
type Table struct {
	Writer *tabwriter.Writer
	Buffer *bytes.Buffer
}

// NewTable :
func NewTable(columns ...string) (*Table, error) {
	var b bytes.Buffer
	minwidth, tabwidth := 1, 1
	padding := 2
	noflag := uint(0)
	w := tabwriter.NewWriter(&b, minwidth, tabwidth, padding, ' ', noflag)
	if _, err := w.Write([]byte(strings.Join(columns, "\t|\t") + "\n")); err != nil {
		return nil, err
	}
	return &Table{Writer: w, Buffer: &b}, nil
}

// AddLine :
func (table *Table) AddLine(cells ...string) error {
	if _, err := table.Writer.Write([]byte(strings.Join(cells, "\t|\t") + "\n")); err != nil {
		return err
	}
	return nil
}

// Lines :
func (table *Table) Lines() ([][]byte, error) {
	if err := table.Writer.Flush(); err != nil {
		return nil, errors.WithStack(err)
	}
	lines := bytes.Split(table.Buffer.Bytes(), []byte("\n"))
	return lines[:len(lines)-1], nil
}

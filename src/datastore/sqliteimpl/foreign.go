package sqliteimpl

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// ForeignKey : foreign key constraint
type ForeignKey struct {
	References string
	ColumnName string
}

func (key ForeignKey) String() string {
	return fmt.Sprintf("FOREIGN KEY(%s) REFERENCES %s", key.ColumnName, key.References)
}

// ForeignKeys :
type ForeignKeys []ForeignKey

func (keys ForeignKeys) String() string {
	parts := []string{}
	for _, key := range keys {
		parts = append(parts, ", "+key.String())
	}
	return strings.Join(parts, "")
}

func (keys *ForeignKeys) gather(base interface{}) error {
	v := reflect.TypeOf(base)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if unicode.IsLower((rune(field.Name[0]))) {
			continue
		}

		if field.Anonymous {
			f := reflect.New(field.Type).Elem().Interface()
			if err := keys.gather(f); err != nil {
				return errors.WithStack(err)
			}
			continue
		}

		foreignTag, ok := field.Tag.Lookup("foreign")
		if !ok {
			continue
		}

		name, ok := columnName(field)
		if !ok {
			continue
		}

		*keys = append(*keys, ForeignKey{
			References: foreignTag,
			ColumnName: name,
		})
	}
	return nil
}

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

func (key ForeignKey) toSQLPart() string {
	return fmt.Sprintf(", FOREIGN KEY(%s) REFERENCES %s", key.ColumnName, key.References)
}

// ForeignKeys :
type ForeignKeys []ForeignKey

func (keys *ForeignKeys) gather(val interface{}) error {
	v := reflect.TypeOf(val)
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

		dbTag, ok := field.Tag.Lookup("db")
		if !ok {
			continue
		}
		columnName := strings.Split(dbTag, ",")[0]

		key := ForeignKey{
			References: foreignTag,
			ColumnName: columnName,
		}
		*keys = append(*keys, key)
	}
	return nil
}

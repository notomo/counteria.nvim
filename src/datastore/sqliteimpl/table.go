package sqliteimpl

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-gorp/gorp"
	"github.com/pkg/errors"
)

// Table :
type Table struct {
	Base interface{}
	Name string
}

var sqlSuffix = regexp.MustCompile(`\)\s*;$`)

// Create :
func (table *Table) Create(dbmap *gorp.DbMap) error {
	checks := Checks{}
	if err := checks.gather(table.Base); err != nil {
		return errors.WithStack(err)
	}

	foreignKeys := ForeignKeys{}
	if err := foreignKeys.gather(table.Base); err != nil {
		return errors.WithStack(err)
	}

	sqlParts := checks.String() + foreignKeys.String()
	ifNotExists := true
	baseSQL := dbmap.AddTableWithName(table.Base, table.Name).SqlForCreate(ifNotExists)
	sql := sqlSuffix.ReplaceAllString(baseSQL, sqlParts+") ;")
	if _, err := dbmap.Exec(sql); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Tables :
type Tables []Table

// Setup :
func (tables Tables) Setup(dbmap *gorp.DbMap) error {
	for _, table := range tables {
		if err := table.Create(dbmap); err != nil {
			return errors.WithStack(err)
		}
	}

	if _, err := dbmap.Exec("PRAGMA foreign_keys=true"); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func columnName(field reflect.StructField) (string, bool) {
	dbTag, ok := field.Tag.Lookup("db")
	if !ok {
		return "", false
	}
	return strings.Split(dbTag, ",")[0], true
}

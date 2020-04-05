package sqliteimpl

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-gorp/gorp"
	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/pkg/errors"
)

// CheckMap :
type CheckMap map[string]Check

func (c CheckMap) gather(val interface{}) error {
	v := reflect.TypeOf(val)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if unicode.IsLower((rune(field.Name[0]))) {
			continue
		}

		if field.Anonymous {
			f := reflect.New(field.Type).Elem().Interface()
			if err := c.gather(f); err != nil {
				return errors.WithStack(err)
			}
			continue
		}

		checkTag, ok := field.Tag.Lookup("check")
		if !ok {
			continue
		}
		fn, ok := checkFuncs[checkTag]
		if !ok {
			return errors.New("invalid check tag: " + checkTag)
		}

		dbTag, ok := field.Tag.Lookup("db")
		if !ok {
			continue
		}
		columnName := strings.Split(dbTag, ",")[0]

		c[columnName] = Check{
			ColumnName: columnName,
			Func:       fn,
		}
	}
	return nil
}

// Check : check constraint
type Check struct {
	Func       func(column string) string
	ColumnName string
}

func (c Check) toSQLPart() string {
	return fmt.Sprintf(", CHECK (%s)", c.Func(c.ColumnName))
}

var checkFuncs = map[string]func(string) string{
	"notEmpty": func(column string) string {
		return fmt.Sprintf(`%s != ""`, column)
	},
	"natural": func(column string) string {
		return fmt.Sprintf(`%s > 0`, column)
	},
	"periodUnit": func(column string) string {
		enums := []string{}
		for _, e := range model.PeriodUnits() {
			enums = append(enums, fmt.Sprintf(`"%s"`, e))
		}
		return fmt.Sprintf(`%s IN (%s)`, column, strings.Join(enums, ", "))
	},
}

var sqlSuffixPattern = regexp.MustCompile(`\)\s*;$`)

func toCreateSQL(table *gorp.TableMap, checks CheckMap, keys ForeignKeys) string {
	checkParts := []string{}
	for _, col := range table.Columns {
		check, ok := checks[col.ColumnName]
		if !ok {
			continue
		}
		checkParts = append(checkParts, check.toSQLPart())
	}
	for _, key := range keys {
		checkParts = append(checkParts, key.toSQLPart())
	}
	checkSQL := strings.Join(checkParts, "")

	ifNotExists := true
	base := table.SqlForCreate(ifNotExists)
	sql := sqlSuffixPattern.ReplaceAllString(base, checkSQL+") ;")
	return sql
}

func createTable(dbmap *gorp.DbMap, base interface{}, name string) error {
	table := dbmap.AddTableWithName(base, name)

	checks := CheckMap{}
	if err := checks.gather(base); err != nil {
		return errors.WithStack(err)
	}

	foreignKeys := ForeignKeys{}
	if err := foreignKeys.gather(base); err != nil {
		return errors.WithStack(err)
	}

	sql := toCreateSQL(table, checks, foreignKeys)
	if _, err := dbmap.Exec(sql); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

package sqliteimpl

import (
	"fmt"
	"log"
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

func gatherChecks(val interface{}, result CheckMap) {
	v := reflect.TypeOf(val)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		log.Println(field.Name)
		if unicode.IsLower((rune(field.Name[0]))) {
			continue
		}

		if field.Anonymous {
			f := reflect.New(field.Type).Elem().Interface()
			gatherChecks(f, result)
			continue
		}

		checkTag, ok := field.Tag.Lookup("check")
		if !ok {
			continue
		}
		if _, ok := checkFuncs[checkTag]; !ok {
			continue
		}

		dbTag, ok := field.Tag.Lookup("db")
		if !ok {
			continue
		}
		columnName := strings.Split(dbTag, ",")[0]

		result[columnName] = Check{
			ColumnName: columnName,
			Name:       checkTag,
		}
	}
}

// Check : check constraint
type Check struct {
	Name       string
	ColumnName string
}

func (c Check) toSQLPart() string {
	return fmt.Sprintf(", CHECK (%s)", checkFuncs[c.Name](c.ColumnName))
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

func toCreateSQL(table *gorp.TableMap, checks CheckMap) string {
	checkParts := []string{}
	for _, col := range table.Columns {
		check, ok := checks[col.ColumnName]
		if !ok {
			continue
		}
		checkParts = append(checkParts, check.toSQLPart())
	}

	checkSQL := strings.Join(checkParts, "")

	ifNotExists := true
	base := table.SqlForCreate(ifNotExists)
	sql := sqlSuffixPattern.ReplaceAllString(base, checkSQL+") ;")
	return sql
}

func createTable(dbmap *gorp.DbMap, table *gorp.TableMap, base interface{}) error {
	checks := CheckMap{}
	gatherChecks(base, checks)

	sql := toCreateSQL(table, checks)
	if _, err := dbmap.Exec(sql); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

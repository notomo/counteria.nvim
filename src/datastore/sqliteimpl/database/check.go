package database

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/pkg/errors"
)

// Check : check constraint
type Check struct {
	Fn         func(column string) string
	ColumnName string

	RawCheck string
}

func (c Check) String() string {
	if c.RawCheck != "" {
		return fmt.Sprintf("CHECK (%s)", c.RawCheck)
	}
	return fmt.Sprintf("CHECK (%s)", c.Fn(c.ColumnName))
}

// Checks :
type Checks []Check

func (checks Checks) String() string {
	parts := []string{}
	for _, check := range checks {
		parts = append(parts, ", "+check.String())
	}
	return strings.Join(parts, "")
}

func (checks *Checks) gather(base interface{}, rawChecks ...string) error {
	v := reflect.TypeOf(base)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if unicode.IsLower((rune(field.Name[0]))) {
			continue
		}

		if field.Anonymous {
			f := reflect.New(field.Type).Elem().Interface()
			if err := checks.gather(f); err != nil {
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

		name, ok := columnName(field)
		if !ok {
			continue
		}

		*checks = append(*checks, Check{
			ColumnName: name,
			Fn:         fn,
		})
	}

	for _, raw := range rawChecks {
		*checks = append(*checks, Check{RawCheck: raw})
	}

	return nil
}

var checkFuncs = map[string]func(string) string{
	"notEmpty": func(column string) string {
		return fmt.Sprintf(`%s != ""`, column)
	},
	"natural": func(column string) string {
		return fmt.Sprintf(`%s > 0`, column)
	},
	"weekday": func(column string) string {
		enums := []string{}
		for _, e := range model.AllWeekdays() {
			enums = append(enums, fmt.Sprintf(`%d`, e))
		}
		return fmt.Sprintf(`%s IN (%s)`, column, strings.Join(enums, ", "))
	},
	"day": func(column string) string {
		return fmt.Sprintf(`1 <= %s AND %s <= 31`, column, column)
	},
	"periodUnit": func(column string) string {
		enums := []string{}
		for _, e := range model.PeriodUnits() {
			enums = append(enums, fmt.Sprintf(`"%s"`, e))
		}
		return fmt.Sprintf(`%s IN (%s)`, column, strings.Join(enums, ", "))
	},
	"taskRuleType": func(column string) string {
		enums := []string{}
		for _, e := range model.TaskRuleTypes() {
			enums = append(enums, fmt.Sprintf(`"%s"`, e))
		}
		return fmt.Sprintf(`%s IN (%s)`, column, strings.Join(enums, ", "))
	},
}

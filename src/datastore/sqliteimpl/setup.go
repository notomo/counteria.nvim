package sqliteimpl

import (
	"github.com/notomo/counteria.nvim/src/datastore/sqliteimpl/database"
	"github.com/notomo/counteria.nvim/src/domain"
	"github.com/pkg/errors"
)

// Setup : tables, dependencies
func Setup(opts ...func(*database.Config)) (*domain.Dep, error) {
	config := &database.Config{}
	for _, opt := range opts {
		opt(config)
	}

	tables := database.Tables{
		{Base: Task{}, Name: "tasks"},
		{Base: DoneTask{}, Name: "done_tasks"},
		{
			Base:      TaskRuleLine{},
			Name:      "task_rule_lines",
			RawChecks: ruleLineChecks,
		},
	}
	dbmap, err := database.Setup(tables, config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &domain.Dep{
		TaskRepository: &TaskRepository{
			Db:    dbmap,
			Rules: &TaskRuleLineRepository{Db: dbmap},
			Dones: &DoneTaskRepository{Db: dbmap},
		},
		TransactionFactory: &TransactionFactory{Db: dbmap},
	}, nil
}

var ruleLineChecks = []string{
	`(
		weekday IS NOT NULL
		AND day IS NULL
		AND month_day IS NULL
		AND date_time IS NULL
		AND rule_date IS NULL
		AND period_number IS NULL
		AND period_unit IS NULL
	) OR
		weekday IS NULL`,
	`(
		weekday IS NULL
		AND day IS NOT NULL
		AND month_day IS NULL
		AND date_time IS NULL
		AND rule_date IS NULL
		AND period_number IS NULL
		AND period_unit IS NULL
	) OR
		day IS NULL`,
	`(
		weekday IS NULL
		AND day IS NULL
		AND month_day IS NOT NULL
		AND date_time IS NULL
		AND rule_date IS NULL
		AND period_number IS NULL
		AND period_unit IS NULL
	) OR
		month_day IS NULL`,
	`(
		weekday IS NULL
		AND day IS NULL
		AND month_day IS NULL
		AND date_time IS NOT NULL
		AND rule_date IS NULL
		AND period_number IS NULL
		AND period_unit IS NULL
	) OR
		date_time IS NULL`,
	`(
		weekday IS NULL
		AND day IS NULL
		AND month_day IS NULL
		AND date_time IS NULL
		AND rule_date IS NOT NULL
		AND period_number IS NULL
		AND period_unit IS NULL
	) OR
		rule_date IS NULL`,
	`(
		weekday IS NULL
		AND day IS NULL
		AND month_day IS NULL
		AND date_time IS NULL
		AND rule_date IS NULL
		AND period_number IS NOT NULL
		AND period_unit IS NOT NULL
	) OR (
		period_number IS NULL
		AND period_unit IS NULL
	)`,
	`NOT (
		weekday IS NULL
		AND day IS NULL
		AND month_day IS NULL
		AND date_time IS NULL
		AND rule_date IS NULL
		AND period_number IS NULL
		AND period_unit IS NULL
	)`,
}

// WithDataPath :
func WithDataPath(path string) func(*database.Config) {
	return func(op *database.Config) {
		op.DataPath = path
	}
}

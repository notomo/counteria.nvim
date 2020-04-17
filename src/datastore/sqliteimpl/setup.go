package sqliteimpl

import (
	"github.com/notomo/counteria.nvim/src/datastore/sqliteimpl/database"
	"github.com/notomo/counteria.nvim/src/domain"
	"github.com/pkg/errors"
)

// Setup : tables, dependencies
func Setup() (*domain.Dep, error) {
	tables := database.Tables{
		{Base: Task{}, Name: "tasks"},
		{Base: DoneTask{}, Name: "done_tasks"},
		{Base: TaskRuleLine{}, Name: "task_rule_lines"},
	}
	dbmap, err := database.Setup(tables)
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

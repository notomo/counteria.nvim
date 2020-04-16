package sqliteimpl

import (
	"time"

	"github.com/go-gorp/gorp"
	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/notomo/counteria.nvim/src/domain/repository"
	"github.com/pkg/errors"
)

// DoneTaskRepository :
type DoneTaskRepository struct {
	Db *gorp.DbMap
}

// Create :
func (repo *DoneTaskRepository) Create(transaction repository.Transaction, task *model.Task, now time.Time) error {
	trans := transaction.(*gorp.Transaction)

	done := DoneTask{
		TaskID:   task.ID(),
		TaskName: task.Name(),
		DoneAt:   now,
	}
	if err := trans.Insert(&done); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Delete :
func (repo *DoneTaskRepository) Delete(transaction repository.Transaction, task *model.Task) error {
	trans := transaction.(*gorp.Transaction)

	dones, err := repo.List(task.ID())
	if err != nil {
		return errors.WithStack(err)
	}
	ds := make([]interface{}, len(dones))
	for i, d := range dones {
		d := d
		ds[i] = &d
	}

	if _, err := trans.Delete(ds...); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// List :
func (repo *DoneTaskRepository) List(taskID int) ([]DoneTask, error) {
	dones := []DoneTask{}
	if _, err := repo.Db.Select(&dones, `
	SELECT *
	FROM done_tasks
	WHERE task_id = ?
	`, taskID); err != nil {
		return nil, errors.WithStack(err)
	}
	return dones, nil
}

var _ model.DoneTaskData = &DoneTask{}

// DoneTask :
type DoneTask struct {
	DoneTaskID int       `db:"id, primarykey, autoincrement"`
	TaskID     int       `db:"task_id, notnull" foreign:"tasks(id)"`
	TaskName   string    `db:"name, notnull" check:"notEmpty"`
	DoneAt     time.Time `db:"at, notnull"`
}

// At :
func (done *DoneTask) At() time.Time {
	return done.DoneAt
}

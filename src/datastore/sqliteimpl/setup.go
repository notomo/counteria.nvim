package sqliteimpl

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/go-gorp/gorp"
	"github.com/notomo/counteria.nvim/src/domain"
	"github.com/pkg/errors"
)

// Setup : tables, dependencies
func Setup() (*domain.Dep, error) {
	dbDirPath := filepath.Join(xdg.DataHome, "counteria")
	if err := os.MkdirAll(dbDirPath, 0770); err != nil {
		return nil, errors.WithStack(err)
	}

	dbPath := filepath.Join(dbDirPath, "default.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	if err := createTable(dbmap, Task{}, "tasks"); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := createTable(dbmap, DoneTask{}, "done_tasks"); err != nil {
		return nil, errors.WithStack(err)
	}
	if _, err := dbmap.Exec("PRAGMA foreign_keys=true"); err != nil {
		return nil, errors.WithStack(err)
	}

	return &domain.Dep{
		TaskRepository:     &TaskRepository{Db: dbmap},
		TransactionFactory: &TransactionFactory{Db: dbmap},
	}, nil
}

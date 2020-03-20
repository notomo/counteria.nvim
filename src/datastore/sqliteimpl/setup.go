package sqliteimpl

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/go-gorp/gorp"
	"github.com/notomo/counteria.nvim/src/domain"
	"github.com/notomo/counteria.nvim/src/domain/model"
)

// Setup : tables, dependencies
func Setup() (*domain.Dep, error) {
	dbDirPath := filepath.Join(xdg.DataHome, "counteria")
	if err := os.MkdirAll(dbDirPath, 0770); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dbDirPath, "default.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(model.Task{}, "tasks").SetKeys(true, "ID")

	if err := dbmap.CreateTablesIfNotExists(); err != nil {
		return nil, err
	}

	return &domain.Dep{
		TaskRepository: &TaskRepository{Db: dbmap},
	}, nil
}

package database

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/go-gorp/gorp"
	"github.com/pkg/errors"
)

// Setup : database file, tables
func Setup(tables Tables) (*gorp.DbMap, error) {
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
	dbmap.ExpandSliceArgs = true

	if err := tables.Setup(dbmap); err != nil {
		return nil, errors.WithStack(err)
	}

	return dbmap, nil
}

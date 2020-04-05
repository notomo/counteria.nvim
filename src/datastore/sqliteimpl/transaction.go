package sqliteimpl

import (
	"github.com/go-gorp/gorp"
	"github.com/notomo/counteria.nvim/src/domain/repository"
	"github.com/pkg/errors"
)

var _ repository.TransactionFactory = &TransactionFactory{}

// TransactionFactory : impl
type TransactionFactory struct {
	Db *gorp.DbMap
}

// Begin :
func (factory *TransactionFactory) Begin() (repository.Transaction, error) {
	trans, err := factory.Db.Begin()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return trans, nil
}

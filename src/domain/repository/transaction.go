package repository

// TransactionFactory :
type TransactionFactory interface {
	Begin() (Transaction, error)
}

// Transaction :
type Transaction interface {
	Commit() error
	Rollback() error
}

package domain

import "github.com/notomo/counteria.nvim/src/domain/repository"

// Dep : dependencies
type Dep struct {
	TaskRepository     repository.TaskRepository
	TransactionFactory repository.TransactionFactory
}

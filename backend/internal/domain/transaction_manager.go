package domain

import "context"

type TransactionManager interface {
	TxBegin() TransactionManager
	TxCommit() error
	TxRollback()
	GetTx() interface{}
	EndTx(error) error
	RecoverTx()
	AssignToContext(parentCtx context.Context) context.Context
}

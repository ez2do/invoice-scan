package repo

import (
	"context"

	"invoice-scan/backend/internal/domain"
	"invoice-scan/backend/pkg/log"

	"gorm.io/gorm"
)

type GormTransactionManager struct {
	db     *gorm.DB
	isDone bool
}

func NewGormTransactionManager(db *gorm.DB) domain.TransactionManager {
	return &GormTransactionManager{db: db}
}

func (tm *GormTransactionManager) TxBegin() domain.TransactionManager {
	return &GormTransactionManager{
		db: tm.db.Begin(),
	}
}

func (tm *GormTransactionManager) TxCommit() (err error) {
	if !tm.isDone {
		err = tm.db.Commit().Error
		tm.isDone = true
	}
	return err
}

func (tm *GormTransactionManager) TxRollback() {
	if !tm.isDone {
		tm.db.Rollback()
		tm.isDone = true
	}
}

func (tm *GormTransactionManager) GetTx() interface{} {
	return tm.db
}

func (tm *GormTransactionManager) EndTx(err error) error {
	if err != nil {
		log.Errorf("transaction: found error, rolling back: %v", err)
		tm.TxRollback()
	} else {
		err = tm.TxCommit()
		if err != nil {
			log.Errorf("transaction: found error when commit, rolling back: %v", err)
			tm.TxRollback()
		}
	}

	return err
}

func (tm *GormTransactionManager) RecoverTx() {
	if p := recover(); p != nil {
		log.Errorf("transaction: found panic, rolling back: %v", p)
		tm.TxRollback()
		panic(p) // Re-panic after rollback
	}
}

type contextKey string

const (
	ContextKeyRepoTx contextKey = "context_key_repo_tx"
)

func (tm *GormTransactionManager) AssignToContext(parentCtx context.Context) context.Context {
	return context.WithValue(parentCtx, ContextKeyRepoTx, tm.GetTx())
}

func getDBFromContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	ctxDB := ctx.Value(ContextKeyRepoTx)
	if ctxDB == nil {
		return db
	}

	txDB, ok := ctxDB.(*gorm.DB)
	if !ok {
		return db
	}

	if txDB != nil {
		return txDB.WithContext(ctx)
	}

	return db
}

package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/avast/retry-go"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var txKey = struct{}{}

type dbTransaction struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) repository.Transaction {
	return &dbTransaction{db: db}
}

func (t *dbTransaction) DoInTx(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	var value interface{}

	txFn := func() error {
		v, err := t.commit(ctx, f)
		if err != nil {
			return err
		}
		value = v
		return nil
	}

	err := retry.Do(
		txFn,
		retry.RetryIf(func(err error) bool {
			return t.isDeadLock(err)
		}),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return time.Duration(n) * time.Second
		}),
		retry.Attempts(4),
		retry.LastErrorOnly(true),
	)

	return value, err
}

func GetTx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(&txKey).(*gorm.DB)
	return tx, ok
}

func (t *dbTransaction) commit(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (value interface{}, err error) {
	tx := t.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, &txKey, tx)

	v, err := f(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return v, nil
}

func (t *dbTransaction) isDeadLock(err error) bool {
	if err == nil {
		return false
	}

	errDeadLock := &mysql.MySQLError{Number: 1213}
	if errors.As(err, &errDeadLock) {
		for err != nil {
			if errors.Is(err, errDeadLock) {
				return true
			}
			err = errors.Unwrap(err)
		}
	}

	return false
}

package repository

import "errors"

var (
	ErrTx       = errors.New("transaction error")
	ErrNotFound = errors.New("record not found")
)

type PreloadCondition struct {
	Query     string
	Condition []any
}

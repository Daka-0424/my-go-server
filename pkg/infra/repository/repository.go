package repository

import (
	"context"

	"gorm.io/gorm"
)

func getTx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(&txKey).(*gorm.DB)
	return tx, ok
}

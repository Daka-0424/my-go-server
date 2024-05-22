package repository

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
)

type IAdmin interface {
	// Exsists checks if the admin with the given email exists in the database.
	Exsists(ctx context.Context, email string) bool

	// Register creates a new admin with the given email, password, and roleType.
	Register(ctx context.Context, email, pass string, roleType entity.AdminRoleType) (*entity.Admin, error)

	// Find returns a list of admins that match the given conditions.
	Find(ctx context.Context, param entity.Admin) ([]entity.Admin, error)

	// Update updates the admin with the given email.
	Update(ctx context.Context, admin *entity.Admin) error

	// Delete deletes the admin with the given email.
	Delete(ctx context.Context, admin *entity.Admin) error

	// CountAll returns the total number of admins in the database.
	CountAll(ctx context.Context) (int64, error)
}

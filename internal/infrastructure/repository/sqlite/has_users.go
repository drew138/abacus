package sqlite

import (
	"context"
)

func (r *Repository) HasUsers(ctx context.Context) (bool, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM users)").Scan(&exists); err != nil {
		return false, domainError(err, "checking for users")
	}
	return exists, nil
}

package rqlite

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) HasUsers(ctx context.Context) (bool, error) {
	qr, err := r.conn.QueryOneContext(ctx, "SELECT EXISTS (SELECT 1 FROM users)")
	if err != nil {
		return false, domainError(err, "checking for users")
	}
	var exists bool
	if qr.Next() {
		if err := qr.Scan(&exists); err != nil {
			return false, fmt.Errorf("reading whether users exist: %w: %w", domain.ErrInternal, err)
		}
	}
	return exists, nil
}

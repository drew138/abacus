package sqlite

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) FindSession(ctx context.Context, tokenHash string) (domain.Session, error) {
	var row sessionRow
	err := r.db.QueryRowContext(ctx, "SELECT token_hash, user_id, expires_at FROM sessions WHERE token_hash = ?", tokenHash).
		Scan(&row.TokenHash, &row.UserID, &row.ExpiresAt)
	if err != nil {
		return domain.Session{}, domainError(err, "selecting a session")
	}
	return row.toDomain(), nil
}

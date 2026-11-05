package sqlite

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	var row findByUsernameRow
	err := r.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE username = ?", username).
		Scan(&row.ID, &row.Username, &row.PasswordHash)
	if err != nil {
		return domain.User{}, domainError(err, "selecting a user by name")
	}
	return row.toDomain(), nil
}

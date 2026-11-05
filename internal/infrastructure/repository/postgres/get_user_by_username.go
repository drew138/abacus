package postgres

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	var row getUserByUsernameRow
	err := r.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE username = $1", username).
		Scan(&row.ID, &row.Username, &row.PasswordHash)
	if err != nil {
		return domain.User{}, domainError(err, "selecting a user")
	}
	return row.toDomain(), nil
}

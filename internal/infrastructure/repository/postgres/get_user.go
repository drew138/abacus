package postgres

import (
	"context"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) GetUser(ctx context.Context, id int64) (domain.User, error) {
	var row getUserRow
	err := r.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE id = $1", id).
		Scan(&row.ID, &row.Username, &row.PasswordHash)
	if err != nil {
		return domain.User{}, domainError(err, "selecting a user")
	}
	user := row.toDomain()
	user.PasswordHash = ""
	return user, nil
}

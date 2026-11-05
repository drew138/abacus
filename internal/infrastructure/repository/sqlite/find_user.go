package sqlite

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

func (r *Repository) FindUser(ctx context.Context, id int64) (domain.User, error) {
	var row findUserRow
	err := r.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE id = ?", id).
		Scan(&row.ID, &row.Username, &row.PasswordHash)
	if err != nil {
		return domain.User{}, domainError(err, fmt.Sprintf("selecting user %d", id))
	}
	user := row.toDomain()
	user.PasswordHash = ""
	return user, nil
}

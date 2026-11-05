package postgres

import (
	"context"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/password"
)

func (r *Repository) CreateFirstUser(ctx context.Context, username, plain string) (domain.User, error) {
	hash, err := password.Hash(plain)
	if err != nil {
		return domain.User{}, fmt.Errorf("hashing the password: %w: %w", domain.ErrInternal, err)
	}
	user := domain.User{Username: username}
	err = r.db.QueryRowContext(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id", username, hash).Scan(&user.ID)
	if err != nil {
		return domain.User{}, domainError(err, "creating the first user")
	}
	return user, nil
}

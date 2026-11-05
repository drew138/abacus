package postgres

import "github.com/drew138/abacus/internal/domain"

type getUserRow struct {
	ID           int64
	Username     string
	PasswordHash string
}

func (r getUserRow) toDomain() domain.User {
	return domain.User{ID: r.ID, Username: r.Username, PasswordHash: r.PasswordHash}
}

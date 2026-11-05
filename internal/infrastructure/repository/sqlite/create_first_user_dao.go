package sqlite

import "github.com/drew138/abacus/internal/domain"

type createFirstUserRow struct {
	ID           int64
	Username     string
	PasswordHash string
}

func (r createFirstUserRow) toDomain() domain.User {
	return domain.User{ID: r.ID, Username: r.Username, PasswordHash: r.PasswordHash}
}

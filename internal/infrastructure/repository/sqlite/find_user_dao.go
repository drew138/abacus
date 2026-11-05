package sqlite

import "github.com/drew138/abacus/internal/domain"

type findUserRow struct {
	ID           int64
	Username     string
	PasswordHash string
}

func (r findUserRow) toDomain() domain.User {
	return domain.User{ID: r.ID, Username: r.Username, PasswordHash: r.PasswordHash}
}

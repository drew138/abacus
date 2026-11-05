package sqlite

import "github.com/drew138/abacus/internal/domain"

type findByUsernameRow struct {
	ID           int64
	Username     string
	PasswordHash string
}

func (r findByUsernameRow) toDomain() domain.User {
	return domain.User{ID: r.ID, Username: r.Username, PasswordHash: r.PasswordHash}
}

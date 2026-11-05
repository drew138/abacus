package postgres

import "context"

func (r *Repository) CountUsers(ctx context.Context) (int, error) {
	var n int
	if err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM users").Scan(&n); err != nil {
		return 0, domainError(err, "counting users")
	}
	return n, nil
}

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	driver "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"github.com/drew138/abacus/internal/domain"
)

const dateLayout = "2006-01-02"

var constraintErrors = map[string]error{
	"users.username": domain.ErrConflict,
	"movements.ref":  domain.ErrConflict,
}

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

const movementColumns = "id, ref, occurred_on, merchant, merchant_key, amount, category, method, note, source, edited, deleted_at"

func scanMovement(scan func(...any) error) (movementRow, error) {
	var r movementRow
	err := scan(&r.ID, &r.Ref, &r.OccurredOn, &r.Merchant, &r.MerchantKey, &r.Amount, &r.Category, &r.Method, &r.Note, &r.Source, &r.Edited, &r.DeletedAt)
	return r, err
}

func (r *Repository) Import(ctx context.Context, movements []domain.Movement) (domain.ImportResult, error) {
	var result domain.ImportResult
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return result, domainError(err, "starting import")
	}
	defer func() { _ = tx.Rollback() }()

	at := now()
	for _, m := range movements {
		row := mapMovementRow(m)
		var id int64
		var edited bool
		var deletedAt sql.NullString
		err := tx.QueryRowContext(ctx, "SELECT id, edited, deleted_at FROM movements WHERE ref = ?", row.Ref).Scan(&id, &edited, &deletedAt)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO movements (ref, occurred_on, merchant, merchant_key, amount, category, method, note, source, edited, created_at, updated_at)
				 VALUES (?, ?, ?, ?, ?, ?, ?, '', ?, 0, ?, ?)`,
				row.Ref, row.OccurredOn, row.Merchant, row.MerchantKey, row.Amount, row.Category, row.Method, row.Source, at, at); err != nil {
				return domain.ImportResult{}, domainError(err, "inserting an imported movement")
			}
			result.Added++
		case err != nil:
			return domain.ImportResult{}, domainError(err, "finding an imported movement")
		case edited || deletedAt.Valid:
			result.Kept++
		default:
			res, err := tx.ExecContext(ctx,
				`UPDATE movements SET occurred_on = ?, merchant = ?, merchant_key = ?, amount = ?, category = ?, method = ?, updated_at = ?
				 WHERE id = ? AND (occurred_on <> ? OR merchant <> ? OR amount <> ? OR category <> ? OR method <> ?)`,
				row.OccurredOn, row.Merchant, row.MerchantKey, row.Amount, row.Category, row.Method, at,
				id, row.OccurredOn, row.Merchant, row.Amount, row.Category, row.Method)
			if err != nil {
				return domain.ImportResult{}, domainError(err, fmt.Sprintf("updating imported movement %d", id))
			}
			if n, _ := res.RowsAffected(); n > 0 {
				result.Updated++
			} else {
				result.Kept++
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return domain.ImportResult{}, domainError(err, "committing import")
	}
	return result, nil
}

func (r *Repository) Months(ctx context.Context) ([]domain.MonthSummary, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT substr(occurred_on, 1, 7) AS month,
		        coalesce(sum(CASE WHEN category = ? THEN 0 ELSE amount END), 0),
		        coalesce(sum(CASE WHEN category = ? THEN amount ELSE 0 END), 0),
		        coalesce(sum(CASE WHEN category = ? THEN 0 ELSE 1 END), 0)
		 FROM movements WHERE deleted_at IS NULL
		 GROUP BY month ORDER BY month`,
		domain.CategorySavings, domain.CategorySavings, domain.CategorySavings)
	if err != nil {
		return nil, domainError(err, "selecting months")
	}
	defer rows.Close()

	var months []domain.MonthSummary
	for rows.Next() {
		var row monthRow
		if err := rows.Scan(&row.Month, &row.Spent, &row.Saved, &row.Movements); err != nil {
			return nil, domainError(err, "reading a month")
		}
		months = append(months, row.toDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, domainError(err, "reading months")
	}
	return months, nil
}

func (r *Repository) ListByMonth(ctx context.Context, month string) ([]domain.Movement, error) {
	start, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, fmt.Errorf("parsing month: %w: %w", domain.ErrInvalid, err)
	}
	rows, err := r.db.QueryContext(ctx,
		"SELECT "+movementColumns+" FROM movements WHERE deleted_at IS NULL AND occurred_on >= ? AND occurred_on < ? ORDER BY occurred_on, id",
		start.Format(dateLayout), start.AddDate(0, 1, 0).Format(dateLayout))
	if err != nil {
		return nil, domainError(err, "selecting movements of a month")
	}
	defer rows.Close()

	movements := []domain.Movement{}
	for rows.Next() {
		row, err := scanMovement(rows.Scan)
		if err != nil {
			return nil, domainError(err, "reading a movement")
		}
		movements = append(movements, row.toDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, domainError(err, "reading movements")
	}
	return movements, nil
}

func (r *Repository) Totals(ctx context.Context) ([]domain.CategoryTotal, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT substr(occurred_on, 1, 7) AS month, category, sum(amount)
		 FROM movements WHERE deleted_at IS NULL
		 GROUP BY month, category ORDER BY month, category`)
	if err != nil {
		return nil, domainError(err, "selecting totals")
	}
	defer rows.Close()

	totals := []domain.CategoryTotal{}
	for rows.Next() {
		var row totalRow
		if err := rows.Scan(&row.Month, &row.Category, &row.Total); err != nil {
			return nil, domainError(err, "reading a total")
		}
		totals = append(totals, row.toDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, domainError(err, "reading totals")
	}
	return totals, nil
}

func (r *Repository) Find(ctx context.Context, id int64) (domain.Movement, error) {
	row, err := scanMovement(r.db.QueryRowContext(ctx, "SELECT "+movementColumns+" FROM movements WHERE id = ?", id).Scan)
	if err != nil {
		return domain.Movement{}, domainError(err, fmt.Sprintf("selecting movement %d", id))
	}
	return row.toDomain(), nil
}

func (r *Repository) Create(ctx context.Context, movement domain.Movement) (domain.Movement, error) {
	row := mapMovementRow(movement)
	at := now()
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO movements (ref, occurred_on, merchant, merchant_key, amount, category, method, note, source, edited, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		row.Ref, row.OccurredOn, row.Merchant, row.MerchantKey, row.Amount, row.Category, row.Method, row.Note, row.Source, row.Edited, at, at).Scan(&row.ID)
	if err != nil {
		return domain.Movement{}, domainError(err, "inserting a movement")
	}
	return row.toDomain(), nil
}

func (r *Repository) Update(ctx context.Context, movement domain.Movement) (domain.Movement, error) {
	row := mapMovementRow(movement)
	res, err := r.db.ExecContext(ctx,
		`UPDATE movements SET occurred_on = ?, merchant = ?, merchant_key = ?, amount = ?, category = ?, method = ?, note = ?, edited = ?, updated_at = ?
		 WHERE id = ?`,
		row.OccurredOn, row.Merchant, row.MerchantKey, row.Amount, row.Category, row.Method, row.Note, row.Edited, now(), row.ID)
	if err != nil {
		return domain.Movement{}, domainError(err, fmt.Sprintf("updating movement %d", row.ID))
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return domain.Movement{}, domain.ErrNotFound
	}
	return r.Find(ctx, row.ID)
}

func (r *Repository) SetDeleted(ctx context.Context, id int64, deleted bool) error {
	deletedAt := sql.NullString{String: now(), Valid: deleted}
	res, err := r.db.ExecContext(ctx, "UPDATE movements SET deleted_at = ?, updated_at = ? WHERE id = ?", deletedAt, now(), id)
	if err != nil {
		return domainError(err, fmt.Sprintf("marking movement %d", id))
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Repository) Recategorise(ctx context.Context, merchantKey, category string) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, domainError(err, "starting recategorisation")
	}
	defer func() { _ = tx.Rollback() }()

	at := now()
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO merchant_rules (merchant_key, category, updated_at) VALUES (?, ?, ?)
		 ON CONFLICT (merchant_key) DO UPDATE SET category = excluded.category, updated_at = excluded.updated_at`,
		merchantKey, category, at); err != nil {
		return 0, domainError(err, "saving a merchant rule")
	}
	res, err := tx.ExecContext(ctx,
		"UPDATE movements SET category = ?, updated_at = ? WHERE merchant_key = ? AND category <> ? AND category <> ?",
		category, at, merchantKey, category, domain.CategorySavings)
	if err != nil {
		return 0, domainError(err, "recategorising movements")
	}
	if err := tx.Commit(); err != nil {
		return 0, domainError(err, "committing recategorisation")
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}

func (r *Repository) MerchantRules(ctx context.Context) (map[string]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT merchant_key, category FROM merchant_rules")
	if err != nil {
		return nil, domainError(err, "selecting merchant rules")
	}
	defer rows.Close()

	rules := map[string]string{}
	for rows.Next() {
		var key, category string
		if err := rows.Scan(&key, &category); err != nil {
			return nil, domainError(err, "reading a merchant rule")
		}
		rules[key] = category
	}
	if err := rows.Err(); err != nil {
		return nil, domainError(err, "reading merchant rules")
	}
	return rules, nil
}

func (r *Repository) List(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT name, parent FROM categories ORDER BY name")
	if err != nil {
		return nil, domainError(err, "selecting categories")
	}
	defer rows.Close()

	categories := []domain.Category{}
	for rows.Next() {
		var row categoryRow
		if err := rows.Scan(&row.Name, &row.Parent); err != nil {
			return nil, domainError(err, "reading a category")
		}
		categories = append(categories, row.toDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, domainError(err, "reading categories")
	}
	return categories, nil
}

func (r *Repository) Ensure(ctx context.Context, categories []domain.Category) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domainError(err, "starting to ensure categories")
	}
	defer func() { _ = tx.Rollback() }()

	at := now()
	for _, c := range categories {
		row := mapCategoryRow(c)
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO categories (name, parent, created_at) VALUES (?, ?, ?) ON CONFLICT (name) DO NOTHING",
			row.Name, row.Parent, at); err != nil {
			return domainError(err, "inserting a category")
		}
	}
	if err := tx.Commit(); err != nil {
		return domainError(err, "committing categories")
	}
	return nil
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	var row userRow
	err := r.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE username = ?", username).
		Scan(&row.ID, &row.Username, &row.PasswordHash)
	if err != nil {
		return domain.User{}, domainError(err, "selecting a user by name")
	}
	return row.toDomain(), nil
}

func (r *Repository) FindUser(ctx context.Context, id int64) (domain.User, error) {
	var row userRow
	err := r.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM users WHERE id = ?", id).
		Scan(&row.ID, &row.Username, &row.PasswordHash)
	if err != nil {
		return domain.User{}, domainError(err, fmt.Sprintf("selecting user %d", id))
	}
	user := row.toDomain()
	user.PasswordHash = ""
	return user, nil
}

func (r *Repository) HasUsers(ctx context.Context) (bool, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM users)").Scan(&exists); err != nil {
		return false, domainError(err, "checking for users")
	}
	return exists, nil
}

func (r *Repository) CreateFirstUser(ctx context.Context, username, passwordHash string) (domain.User, error) {
	var row userRow
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (username, password_hash, created_at)
		 SELECT ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM users)
		 RETURNING id, username`,
		username, passwordHash, now()).Scan(&row.ID, &row.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrRegistrationClosed
	}
	if err != nil {
		return domain.User{}, domainError(err, "creating the first user")
	}
	return row.toDomain(), nil
}

func (r *Repository) CreateSession(ctx context.Context, session domain.Session) error {
	row := mapSessionRow(session)
	if _, err := r.db.ExecContext(ctx, "INSERT INTO sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)",
		row.TokenHash, row.UserID, row.ExpiresAt); err != nil {
		return domainError(err, fmt.Sprintf("inserting a session for user %d", row.UserID))
	}
	return nil
}

func (r *Repository) FindSession(ctx context.Context, tokenHash string) (domain.Session, error) {
	var row sessionRow
	err := r.db.QueryRowContext(ctx, "SELECT token_hash, user_id, expires_at FROM sessions WHERE token_hash = ?", tokenHash).
		Scan(&row.TokenHash, &row.UserID, &row.ExpiresAt)
	if err != nil {
		return domain.Session{}, domainError(err, "selecting a session")
	}
	return row.toDomain(), nil
}

func (r *Repository) DeleteSession(ctx context.Context, tokenHash string) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM sessions WHERE token_hash = ?", tokenHash); err != nil {
		return domainError(err, "deleting a session")
	}
	return nil
}

func domainError(err error, doing string) error {
	var sqliteErr *driver.Error
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.ErrNotFound
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrUnavailable, err)
	case errors.As(err, &sqliteErr):
		return sqliteError(sqliteErr, doing)
	default:
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrInternal, err)
	}
}

func sqliteError(err *driver.Error, doing string) error {
	switch err.Code() {
	case sqlite3.SQLITE_CONSTRAINT_UNIQUE, sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
		if mapped, ok := constraintErrors[uniqueTarget(err.Error())]; ok {
			return fmt.Errorf("%s: %w", doing, mapped)
		}
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrConflict, err)
	case sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY, sqlite3.SQLITE_CONSTRAINT_NOTNULL, sqlite3.SQLITE_CONSTRAINT_CHECK:
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrInvalid, err)
	}
	switch err.Code() & 0xff {
	case sqlite3.SQLITE_BUSY, sqlite3.SQLITE_LOCKED, sqlite3.SQLITE_CANTOPEN, sqlite3.SQLITE_IOERR, sqlite3.SQLITE_FULL:
		return fmt.Errorf("%s: %w: %w", doing, domain.ErrUnavailable, err)
	}
	return fmt.Errorf("%s: %w: %w", doing, domain.ErrInternal, err)
}

func uniqueTarget(message string) string {
	_, after, ok := strings.Cut(message, "constraint failed: ")
	if !ok {
		return ""
	}
	target, _, _ := strings.Cut(after, " ")
	return strings.TrimSuffix(target, ",")
}

package sqlite

import (
	"database/sql"
	"time"

	"github.com/drew138/abacus/internal/domain"
)

type movementRow struct {
	ID          int64
	Ref         sql.NullString
	OccurredOn  string
	Merchant    string
	MerchantKey string
	Amount      int64
	Category    string
	Method      string
	Note        string
	Source      string
	Edited      bool
	DeletedAt   sql.NullString
}

func (r movementRow) toDomain() domain.Movement {
	date, _ := time.Parse(dateLayout, r.OccurredOn)
	return domain.Movement{
		ID:       r.ID,
		Ref:      r.Ref.String,
		Date:     date,
		Merchant: r.Merchant,
		Amount:   r.Amount,
		Category: r.Category,
		Method:   domain.Method(r.Method),
		Note:     r.Note,
		Source:   domain.Source(r.Source),
		Edited:   r.Edited,
		Deleted:  r.DeletedAt.Valid,
	}
}

type categoryRow struct {
	Name   string
	Parent sql.NullString
}

func (r categoryRow) toDomain() domain.Category {
	return domain.Category{Name: r.Name, Parent: r.Parent.String}
}

type monthRow struct {
	Month     string
	Spent     int64
	Saved     int64
	Movements int
}

func (r monthRow) toDomain() domain.MonthSummary {
	return domain.MonthSummary{Month: r.Month, Spent: r.Spent, Saved: r.Saved, Movements: r.Movements}
}

type totalRow struct {
	Month    string
	Category string
	Total    int64
}

func (r totalRow) toDomain() domain.CategoryTotal {
	return domain.CategoryTotal{Month: r.Month, Category: r.Category, Total: r.Total}
}

type userRow struct {
	ID           int64
	Username     string
	PasswordHash string
}

func (r userRow) toDomain() domain.User {
	return domain.User{ID: r.ID, Username: r.Username, PasswordHash: r.PasswordHash}
}

type sessionRow struct {
	TokenHash string
	UserID    int64
	ExpiresAt string
}

func (r sessionRow) toDomain() domain.Session {
	expires, _ := time.Parse(time.RFC3339, r.ExpiresAt)
	return domain.Session{TokenHash: r.TokenHash, UserID: r.UserID, ExpiresAt: expires}
}

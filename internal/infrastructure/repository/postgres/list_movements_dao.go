package postgres

import (
	"database/sql"
	"time"

	"github.com/drew138/abacus/internal/domain"
)

type listMovementsRow struct {
	ID         int64
	Ref        sql.NullString
	OccurredOn time.Time
	Merchant   string
	Amount     int64
	Category   string
	Method     string
	Note       string
	Source     string
	Edited     bool
	DeletedAt  sql.NullTime
}

func (r listMovementsRow) toDomain() domain.Movement {
	return domain.Movement{
		ID:       r.ID,
		Ref:      r.Ref.String,
		Date:     time.Date(r.OccurredOn.Year(), r.OccurredOn.Month(), r.OccurredOn.Day(), 0, 0, 0, 0, time.UTC),
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

package rqlite

import (
	"time"

	"github.com/rqlite/gorqlite"

	"github.com/drew138/abacus/internal/domain"
)

type listByMonthMovementRow struct {
	ID         int64
	Ref        gorqlite.NullString
	OccurredOn string
	Merchant   string
	Amount     int64
	Category   string
	Method     string
	Note       string
	Source     string
	Edited     bool
	DeletedAt  gorqlite.NullString
}

func (r listByMonthMovementRow) toDomain() domain.Movement {
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

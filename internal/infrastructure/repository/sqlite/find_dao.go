package sqlite

import (
	"database/sql"
	"time"

	"github.com/drew138/abacus/internal/domain"
)

type findMovementRow struct {
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

func (r findMovementRow) toDomain() domain.Movement {
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

func scanFindMovement(scan func(...any) error) (findMovementRow, error) {
	var r findMovementRow
	err := scan(&r.ID, &r.Ref, &r.OccurredOn, &r.Merchant, &r.MerchantKey, &r.Amount, &r.Category, &r.Method, &r.Note, &r.Source, &r.Edited, &r.DeletedAt)
	return r, err
}

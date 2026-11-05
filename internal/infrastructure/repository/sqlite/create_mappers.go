package sqlite

import (
	"database/sql"

	"github.com/drew138/abacus/internal/domain"
)

func mapCreateMovementRow(m domain.Movement) createMovementRow {
	return createMovementRow{
		ID:          m.ID,
		Ref:         sql.NullString{String: m.Ref, Valid: m.Ref != ""},
		OccurredOn:  m.Date.Format(dateLayout),
		Merchant:    m.Merchant,
		MerchantKey: merchantKey(m.Merchant),
		Amount:      m.Amount,
		Category:    m.Category,
		Method:      string(m.Method),
		Note:        m.Note,
		Source:      string(m.Source),
		Edited:      m.Edited,
	}
}

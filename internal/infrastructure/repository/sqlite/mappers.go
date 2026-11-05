package sqlite

import (
	"database/sql"
	"time"

	"github.com/drew138/abacus/internal/domain"
)

func mapMovementRow(m domain.Movement) movementRow {
	return movementRow{
		ID:          m.ID,
		Ref:         sql.NullString{String: m.Ref, Valid: m.Ref != ""},
		OccurredOn:  m.Date.Format(dateLayout),
		Merchant:    m.Merchant,
		MerchantKey: m.MerchantKey(),
		Amount:      m.Amount,
		Category:    m.Category,
		Method:      string(m.Method),
		Note:        m.Note,
		Source:      string(m.Source),
		Edited:      m.Edited,
	}
}

func mapCategoryRow(c domain.Category) categoryRow {
	return categoryRow{Name: c.Name, Parent: sql.NullString{String: c.Parent, Valid: c.Parent != ""}}
}

func mapSessionRow(s domain.Session) sessionRow {
	return sessionRow{TokenHash: s.TokenHash, UserID: s.UserID, ExpiresAt: s.ExpiresAt.UTC().Format(time.RFC3339)}
}

package rqlite

import "github.com/drew138/abacus/internal/domain"

func mapUpdateMovementRow(m domain.Movement) updateMovementRow {
	return updateMovementRow{
		ID:         m.ID,
		OccurredOn: m.Date.Format(dateLayout),
		Merchant:   m.Merchant,
		Amount:     m.Amount,
		Category:   m.Category,
		Method:     string(m.Method),
		Note:       m.Note,
		Edited:     m.Edited,
	}
}

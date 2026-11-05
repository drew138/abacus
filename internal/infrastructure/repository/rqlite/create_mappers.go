package rqlite

import "github.com/drew138/abacus/internal/domain"

func mapCreateMovementRow(m domain.Movement) createMovementRow {
	var ref any
	if m.Ref != "" {
		ref = m.Ref
	}
	return createMovementRow{
		Ref:         ref,
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

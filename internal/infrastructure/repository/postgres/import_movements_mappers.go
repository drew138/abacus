package postgres

import "github.com/drew138/abacus/internal/domain"

func mapImportMovementsRow(m domain.Movement, rules map[string]string) importMovementsRow {
	row := importMovementsRow{
		Ref:         m.Ref,
		OccurredOn:  m.Date.Format(dateLayout),
		Merchant:    m.Merchant,
		MerchantKey: merchantKey(m.Merchant),
		Amount:      m.Amount,
		Category:    m.Category,
		Method:      string(m.Method),
		Source:      string(m.Source),
	}
	if rule, ok := rules[row.MerchantKey]; ok && row.Category != domain.CategorySavings {
		row.Category = rule
	}
	return row
}

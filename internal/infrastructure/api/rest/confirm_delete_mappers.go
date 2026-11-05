package rest

import (
	"fmt"
	"unicode"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/pesos"
)

func mapConfirmDeleteView(m domain.Movement) confirmDeleteView {
	initial := "?"
	for _, r := range m.Merchant {
		initial = string(unicode.ToUpper(r))
		break
	}
	return confirmDeleteView{
		ID:        m.ID,
		Initial:   initial,
		Color:     categoryColor(m.Category),
		Merchant:  m.Merchant,
		Day:       fmt.Sprintf("%d %s", m.Date.Day(), m.Date.Format("Jan")),
		Category:  m.Category,
		IsSavings: m.IsSavings(),
		Method:    methodLabels[m.Method],
		Manual:    m.Source == domain.SourceManual,
		Edited:    m.Edited,
		Note:      m.Note,
		Amount:    pesos.Format(m.Amount),
		Confirm:   true,
		Month:     m.Month(),
	}
}

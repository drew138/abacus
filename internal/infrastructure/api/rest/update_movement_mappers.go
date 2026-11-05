package rest

import (
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/drew138/abacus/internal/domain"
	"github.com/drew138/abacus/pkg/pesos"
)

func mapUpdateMovementFormView(m domain.Movement, categories []domain.Category, today time.Time, errs []string) updateMovementFormView {
	view := updateMovementFormView{
		ID:            m.ID,
		Title:         "Edit movement",
		Merchant:      m.Merchant,
		Date:          m.Date.Format("2006-01-02"),
		MaxDate:       today.Format("2006-01-02"),
		Note:          m.Note,
		FromStatement: m.Source == domain.SourceStatement,
		Errors:        errs,
	}
	if m.Amount > 0 {
		view.Amount = strings.TrimPrefix(pesos.Format(m.Amount), "$")
	}
	for _, method := range []domain.Method{domain.MethodCredit, domain.MethodDebit} {
		view.Methods = append(view.Methods, updateMovementFormOption{Value: string(method), Label: methodLabels[method], Selected: method == m.Method})
	}
	sorted := append([]domain.Category(nil), categories...)
	sort.Slice(sorted, func(i, j int) bool {
		ti, tj := categoryTier(sorted[i].Name, sorted[i].Parent), categoryTier(sorted[j].Name, sorted[j].Parent)
		if ti != tj {
			return ti < tj
		}
		return sorted[i].Name < sorted[j].Name
	})
	for _, c := range sorted {
		option := updateMovementFormOption{Value: c.Name, Label: c.Name, Color: categoryColor(c.Name), Selected: c.Name == m.Category}
		last := len(view.Categories) - 1
		if last >= 0 && view.Categories[last].Label == c.Parent {
			view.Categories[last].Options = append(view.Categories[last].Options, option)
			continue
		}
		view.Categories = append(view.Categories, updateMovementFormOptionGroup{Label: c.Parent, Options: []updateMovementFormOption{option}})
	}
	return view
}

func mapUpdateMovementRequest(values url.Values) UpdateMovementRequest {
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, values.Get("amount"))
	amount, _ := strconv.ParseInt(digits, 10, 64)
	return UpdateMovementRequest{
		Merchant:        strings.TrimSpace(values.Get("merchant")),
		Amount:          amount,
		Date:            values.Get("date"),
		Category:        values.Get("category"),
		Method:          values.Get("method"),
		Note:            strings.TrimSpace(values.Get("note")),
		ApplyToMerchant: values.Get("apply_to_merchant") == "on",
	}
}

func mapUpdateMovementErrors(err error) []string {
	messages := map[string]string{
		"merchant": "Enter who the money went to, up to 120 characters.",
		"amount":   "Enter the amount in pesos, for example 25.000.",
		"date":     "Pick the date it happened.",
		"category": "Pick a category from the list.",
		"method":   "Pick how it was paid.",
		"note":     "Keep the note under 200 characters.",
	}
	var failed validator.ValidationErrors
	if errors.As(err, &failed) {
		seen := map[string]bool{}
		var out []string
		for _, fe := range failed {
			if msg, ok := messages[fe.Field()]; ok && !seen[msg] {
				seen[msg] = true
				out = append(out, msg)
			}
		}
		return out
	}
	return []string{messages["category"]}
}

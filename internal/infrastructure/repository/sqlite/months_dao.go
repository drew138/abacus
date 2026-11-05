package sqlite

import "github.com/drew138/abacus/internal/domain"

type monthRow struct {
	Month     string
	Spent     int64
	Saved     int64
	Movements int
}

func (r monthRow) toDomain() domain.MonthSummary {
	return domain.MonthSummary{Month: r.Month, Spent: r.Spent, Saved: r.Saved, Movements: r.Movements}
}

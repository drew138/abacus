package rqlite

import "github.com/drew138/abacus/internal/domain"

type monthsRow struct {
	Month     string
	Spent     int64
	Saved     int64
	Movements int64
}

func (r monthsRow) toDomain() domain.MonthSummary {
	return domain.MonthSummary{Month: r.Month, Spent: r.Spent, Saved: r.Saved, Movements: int(r.Movements)}
}

package domain

import "time"

const CategorySavings = "Savings"

type Movement struct {
	ID       int64
	Ref      string
	Date     time.Time
	Merchant string
	Amount   int64
	Category string
	Method   Method
	Note     string
	Source   Source
	Edited   bool
	Deleted  bool
}

func (m Movement) Month() string {
	return m.Date.Format("2006-01")
}

func (m Movement) IsSavings() bool {
	return m.Category == CategorySavings
}

type Category struct {
	Name   string
	Parent string
}

func (c Category) Top() string {
	if c.Parent != "" {
		return c.Parent
	}
	return c.Name
}

type MonthSummary struct {
	Month     string
	Spent     int64
	Saved     int64
	Movements int
}

type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

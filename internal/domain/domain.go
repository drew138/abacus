package domain

import (
	"strings"
	"time"
)

const (
	CategorySavings = "Savings"
	CategoryOther   = "Other"
	GroupFood       = "Food and eating out"
)

func CategoryTier(name, parent string) int {
	switch {
	case name == CategorySavings:
		return 0
	case name == GroupFood || parent == GroupFood:
		return 1
	case name == CategoryOther:
		return 3
	}
	return 2
}

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

func (m Movement) MerchantKey() string {
	return MerchantKey(m.Merchant)
}

func MerchantKey(merchant string) string {
	return strings.ToLower(strings.Join(strings.Fields(merchant), " "))
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

func (c Category) Tier() int {
	return CategoryTier(c.Name, c.Parent)
}

type MonthSummary struct {
	Month     string
	Spent     int64
	Saved     int64
	Movements int
}

type CategoryTotal struct {
	Month    string
	Category string
	Total    int64
}

type ImportResult struct {
	Added   int
	Updated int
	Kept    int
}

type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

type Session struct {
	TokenHash string
	UserID    int64
	ExpiresAt time.Time
}

func (s Session) Expired(now time.Time) bool {
	return !now.Before(s.ExpiresAt)
}

type ExpenseFilter struct {
	Category string
	Search   string
}

type CategoryAmount struct {
	Name        string
	Total       int64
	Previous    int64
	HasPrevious bool
}

type GroupTotal struct {
	CategoryAmount
	Subcategories []CategoryAmount
}

type ExpenseGroup struct {
	Category  string
	Parent    string
	Total     int64
	Movements []Movement
}

type MonthReport struct {
	Month       string
	Summary     MonthSummary
	Previous    MonthSummary
	HasPrevious bool
	Groups      []GroupTotal
	Savings     []Movement
	Expenses    []ExpenseGroup
	Filter      ExpenseFilter
}

type TrendRow struct {
	Name          string
	Values        []int64
	Subcategories []TrendRow
}

type Trends struct {
	Months []MonthSummary
	Rows   []TrendRow
}

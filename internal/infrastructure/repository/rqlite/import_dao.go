package rqlite

type importExistingRow struct {
	Ref       string
	Edited    bool
	DeletedAt string
}

type importMovementRow struct {
	Ref         string
	OccurredOn  string
	Merchant    string
	MerchantKey string
	Amount      int64
	Category    string
	Method      string
	Source      string
}

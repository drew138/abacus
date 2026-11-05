package rqlite

type createMovementRow struct {
	Ref         any
	OccurredOn  string
	Merchant    string
	MerchantKey string
	Amount      int64
	Category    string
	Method      string
	Note        string
	Source      string
	Edited      bool
}

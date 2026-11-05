package postgres

type importMovementsRow struct {
	Ref         string
	OccurredOn  string
	Merchant    string
	MerchantKey string
	Amount      int64
	Category    string
	Method      string
	Source      string
}

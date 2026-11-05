package domain

type Method string

const (
	MethodCredit Method = "credit"
	MethodDebit  Method = "debit"
)

func (m Method) Valid() bool {
	switch m {
	case MethodCredit, MethodDebit:
		return true
	}
	return false
}

type Source string

const (
	SourceStatement Source = "statement"
	SourceManual    Source = "manual"
)

func (s Source) Valid() bool {
	switch s {
	case SourceStatement, SourceManual:
		return true
	}
	return false
}

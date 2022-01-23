package filter

type Comparator uint8

const (
	NOOP = iota // implementors should ignore this filter
	LT
	GT
	LE
	GE
	EQ
)

//Filter is the filter interface
type Filter interface {
	Evaluate() (field string, value interface{}, cmp Comparator)
}

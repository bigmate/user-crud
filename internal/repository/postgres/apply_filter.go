package postgres

import (
	"user-crud/pkg/filter"

	"github.com/Masterminds/squirrel"
)

func applyFilter(qb squirrel.SelectBuilder, f filter.Filter) squirrel.SelectBuilder {
	field, value, comparator := f.Evaluate()
	var op squirrel.Sqlizer
	switch comparator {
	case filter.EQ:
		op = squirrel.Eq{field: value}
	case filter.GT:
		op = squirrel.Gt{field: value}
	case filter.GE:
		op = squirrel.GtOrEq{field: value}
	case filter.LE:
		op = squirrel.LtOrEq{field: value}
	case filter.LT:
		op = squirrel.Lt{field: value}
	default:
		return qb
	}
	return qb.Where(op)
}

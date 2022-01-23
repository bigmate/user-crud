package postgres

import (
	"github.com/Masterminds/squirrel"
)

func pq() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

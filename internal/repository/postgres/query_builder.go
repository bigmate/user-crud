package postgres

import (
	"embed"

	"github.com/Masterminds/squirrel"
)

//go:embed migrations
var Migrations embed.FS

func pq() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

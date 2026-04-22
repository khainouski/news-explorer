// Package source holds the sources-table repository - one file per action.
package source

import "github.com/khainouski/news-explorer/pkg/postgres"

type Repo struct {
	pool *postgres.Pool
}

func NewRepo(pool *postgres.Pool) *Repo {
	return &Repo{pool: pool}
}

// Postgres error codes - see https://www.postgresql.org/docs/current/errcodes-appendix.html.
const (
	pgErrUniqueViolation     = "23505"
	pgErrForeignKeyViolation = "23503"
)

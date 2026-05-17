// Package tag holds the tags-table repository.
package tag

import "github.com/khainouski/news-explorer/pkg/postgres"

type Repo struct {
	pool *postgres.Pool
}

func NewRepo(pool *postgres.Pool) *Repo {
	return &Repo{pool: pool}
}

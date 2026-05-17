// Package article holds the articles-table repository.
package article

import "github.com/khainouski/news-explorer/pkg/postgres"

type Repo struct {
	pool *postgres.Pool
}

func NewRepo(pool *postgres.Pool) *Repo {
	return &Repo{pool: pool}
}

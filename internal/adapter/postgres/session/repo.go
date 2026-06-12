// Package session holds the sessions-table repository.
package session

import "github.com/khainouski/news-explorer/pkg/postgres"

type Repo struct {
	pool *postgres.Pool
}

func NewRepo(pool *postgres.Pool) *Repo {
	return &Repo{pool: pool}
}

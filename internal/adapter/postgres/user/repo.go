// Package user holds the users-table repository. No create.go - there's no signup flow, the one
// account is seeded directly (see migration/postgres/20260507171000_seed.up.sql).
package user

import "github.com/khainouski/news-explorer/pkg/postgres"

type Repo struct {
	pool *postgres.Pool
}

func NewRepo(pool *postgres.Pool) *Repo {
	return &Repo{pool: pool}
}

package article

import "context"

// UnreadCountsBySource returns the count of unread articles per source, across the whole table.
func (u *UseCase) UnreadCountsBySource(ctx context.Context) (map[string]int, error) {
	return u.postgres.UnreadCountsBySource(ctx)
}

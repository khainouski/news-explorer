package source

import "context"

// MarkRead marks every article from a source as read - triggered by clicking the source in the
// home page sidebar (see internal/controller/http/article).
func (u *UseCase) MarkRead(ctx context.Context, sourceID string) error {
	return u.postgres.MarkRead(ctx, sourceID)
}

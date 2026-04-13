package source

import "context"

// Delete removes a source, returning its name (e.g. for a confirmation toast).
func (u *UseCase) Delete(ctx context.Context, id string) (string, error) {
	return u.postgres.Delete(ctx, id)
}

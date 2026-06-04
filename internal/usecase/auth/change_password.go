package auth

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/khainouski/news-explorer/internal/domain"
)

// ChangePassword verifies currentPassword against currentPasswordHash (the caller already has
// the logged-in user's hash from the session, see internal/controller/http/auth) before setting
// newPassword.
func (u *UseCase) ChangePassword(ctx context.Context, userID int, currentPasswordHash, currentPassword, newPassword string) error {
	if bcrypt.CompareHashAndPassword([]byte(currentPasswordHash), []byte(currentPassword)) != nil {
		return domain.ErrInvalidCredentials
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	if err = u.users.UpdatePassword(ctx, userID, string(newHash)); err != nil {
		return fmt.Errorf("users.UpdatePassword: %w", err)
	}

	return nil
}

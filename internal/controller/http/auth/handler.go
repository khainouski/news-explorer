// Package auth handles sign-in/out and the change-password page.
package auth

import usecaseauth "github.com/khainouski/news-explorer/internal/usecase/auth"

type Handler struct {
	auth *usecaseauth.UseCase
}

func NewHandler(auth *usecaseauth.UseCase) *Handler {
	return &Handler{auth: auth}
}

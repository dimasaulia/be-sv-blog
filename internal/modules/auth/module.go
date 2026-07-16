package auth

import (
	"net/http"

	"github.com/sv-blog/internal/modules/auth/controllers"
	"github.com/sv-blog/internal/shared/middleware"
)

type AuthModuleImpl struct {
	AuthController controllers.AuthController
	auth           *middleware.Authenticator
}

func NewAuthModule(authController controllers.AuthController, auth *middleware.Authenticator) *AuthModuleImpl {
	return &AuthModuleImpl{
		AuthController: authController,
		auth:           auth,
	}
}

func (m *AuthModuleImpl) Name() string {
	return "auth"
}

func (m *AuthModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/auth/me", m.auth.RequireAuthenticated(m.AuthController.CurrentUser))
}

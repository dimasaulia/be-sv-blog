package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/sv-blog/internal/platform/authclient"
	"github.com/sv-blog/internal/platform/config"
	"github.com/sv-blog/internal/platform/logger"
	"github.com/sv-blog/internal/shared/requestctx"
	"github.com/sv-blog/internal/shared/response"
)

var (
	errMissingToken = errors.New("middleware: missing bearer token")
	errForbiddenApp = errors.New("middleware: user has no access to this application")
)

type Authenticator struct {
	client   *authclient.Client
	response *response.Sender
	log      *logger.LayerLogger
	appCode  string
}

func NewAuthenticator(cfg config.Config, client *authclient.Client, sender *response.Sender, appLogger *logger.Logger) *Authenticator {
	return &Authenticator{
		client:   client,
		response: sender,
		log:      appLogger.Layer("middleware.authz"),
		appCode:  cfg.Auth.AppCode,
	}
}

func (a *Authenticator) RequireAuthenticated(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ctx, err := a.authenticate(r)
		if err != nil {
			a.respondAuthError(w, r, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Authenticator) LoadUserContext(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ctx, err := a.tryAuthenticate(r)
		if err != nil {
			a.log.Warn(r.Context(), "load_user_context.failed", "error", err.Error())
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Authenticator) RequirePermission(permission string) func(http.HandlerFunc) http.Handler {
	return func(next http.HandlerFunc) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ctx, err := a.authenticate(r)
			if err != nil {
				a.respondAuthError(w, r, err)
				return
			}

			allowed, err := a.client.Check(ctx, token, a.appCode, permission)
			if err != nil {
				a.log.Warn(ctx, "permission.check_failed", "error", err.Error(), "app", a.appCode, "permission", permission)
				a.response.Error(w, r, http.StatusForbidden, "auth.forbidden", nil)
				return
			}
			if !allowed {
				a.response.Error(w, r, http.StatusForbidden, "auth.forbidden", nil)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (a *Authenticator) authenticate(r *http.Request) (string, context.Context, error) {
	token := bearerToken(r)
	if token == "" {
		return "", nil, errMissingToken
	}

	data, err := a.client.Me(r.Context(), token)
	if err != nil {
		return "", nil, err
	}

	if a.appCode != "" && !hasAppAccess(data.Apps, a.appCode) {
		return "", nil, errForbiddenApp
	}

	ctx := requestctx.WithAuthUser(r.Context(), requestctx.AuthUser{
		ID:          data.User.ID,
		Email:       data.User.Email,
		Username:    data.User.Username,
		DisplayName: data.User.DisplayName,
	})
	return token, ctx, nil
}

func (a *Authenticator) tryAuthenticate(r *http.Request) (string, context.Context, error) {
	token := bearerToken(r)
	if token == "" {
		return "", r.Context(), errMissingToken
	}

	data, err := a.client.Me(r.Context(), token)
	if err != nil {
		return "", r.Context(), err
	}

	if a.appCode != "" && !hasAppAccess(data.Apps, a.appCode) {
		return "", r.Context(), errForbiddenApp
	}

	ctx := requestctx.WithAuthUser(r.Context(), requestctx.AuthUser{
		ID:          data.User.ID,
		Email:       data.User.Email,
		Username:    data.User.Username,
		DisplayName: data.User.DisplayName,
	})
	return token, ctx, nil
}

func (a *Authenticator) respondAuthError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, errForbiddenApp) {
		a.response.Error(w, r, http.StatusForbidden, "auth.forbidden", nil)
		return
	}

	if !errors.Is(err, errMissingToken) && !errors.Is(err, authclient.ErrUnauthorized) {
		a.log.Warn(r.Context(), "authenticate.failed", "error", err.Error())
	}
	a.response.Error(w, r, http.StatusUnauthorized, "auth.unauthorized", nil)
}

func hasAppAccess(apps []authclient.App, appCode string) bool {
	for _, app := range apps {
		if strings.EqualFold(app.Code, appCode) {
			return true
		}
	}

	return false
}

func bearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

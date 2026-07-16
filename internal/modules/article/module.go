package article

import (
	"net/http"

	"github.com/sv-blog/internal/modules/article/controllers"
	"github.com/sv-blog/internal/shared/middleware"
	"github.com/sv-blog/internal/shared/permissions"
)

type ArticleModuleImpl struct {
	ArticleController controllers.ArticleController
	auth              *middleware.Authenticator
}

func NewArticleModule(articleController controllers.ArticleController, auth *middleware.Authenticator) *ArticleModuleImpl {
	return &ArticleModuleImpl{
		ArticleController: articleController,
		auth:              auth,
	}
}

func (m *ArticleModuleImpl) Name() string {
	return "article"
}

func (m *ArticleModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /feed/{limit}/{offset}", m.ArticleController.Feed)
	mux.Handle("POST /article/{$}", m.protect(permissions.ArticleWrite, m.ArticleController.Create))
	mux.Handle("GET /article/{limit}/{offset}", m.protect(permissions.ArticleRead, m.ArticleController.Find))
	mux.Handle("GET /article/{id}", m.auth.LoadUserContext(m.ArticleController.FindByID))
	mux.Handle("PUT /article/{id}", m.protect(permissions.ArticleUpdate, m.ArticleController.Update))
	mux.Handle("PATCH /article/{id}", m.protect(permissions.ArticleUpdate, m.ArticleController.Update))
	mux.Handle("DELETE /article/{id}", m.protect(permissions.ArticleDelete, m.ArticleController.Delete))
}

func (m *ArticleModuleImpl) protect(permission string, handler http.HandlerFunc) http.Handler {
	return m.auth.RequirePermission(permission)(handler)
}

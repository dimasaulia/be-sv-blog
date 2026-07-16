//go:build wireinject
// +build wireinject

package app

import (
	"context"

	"github.com/google/wire"
	"github.com/sv-blog/internal/modules/article"
	articleControllers "github.com/sv-blog/internal/modules/article/controllers"
	articleRepositories "github.com/sv-blog/internal/modules/article/repositories"
	articleServices "github.com/sv-blog/internal/modules/article/services"
	"github.com/sv-blog/internal/modules/auth"
	authControllers "github.com/sv-blog/internal/modules/auth/controllers"
	"github.com/sv-blog/internal/modules/health"
	"github.com/sv-blog/internal/modules/health/controllers"
	"github.com/sv-blog/internal/modules/health/repositories"
	"github.com/sv-blog/internal/modules/health/services"
	"github.com/sv-blog/internal/modules/releasenotes"
	releaseNoteControllers "github.com/sv-blog/internal/modules/releasenotes/controllers"
	releaseNoteRepositories "github.com/sv-blog/internal/modules/releasenotes/repositories"
	releaseNoteServices "github.com/sv-blog/internal/modules/releasenotes/services"
	"github.com/sv-blog/internal/platform/authclient"
	"github.com/sv-blog/internal/platform/config"
	"github.com/sv-blog/internal/platform/database"
	"github.com/sv-blog/internal/platform/i18n"
	"github.com/sv-blog/internal/platform/redis"
	"github.com/sv-blog/internal/shared/middleware"
	"github.com/sv-blog/internal/shared/response"
)

func Initialize(ctx context.Context) (*App, error) {
	wire.Build(
		config.Load,
		ProvideLogger,
		database.New,
		redis.New,
		i18n.NewTranslator,
		response.NewSender,
		authclient.NewClient,
		middleware.NewAuthenticator,
		repositories.NewHealthRepository,
		services.NewHealthService,
		controllers.NewHealthController,
		health.NewHealthModule,
		releaseNoteRepositories.NewReleaseNoteRepository,
		releaseNoteServices.NewReleaseNoteService,
		releaseNoteControllers.NewReleaseNoteController,
		releasenotes.NewReleaseNoteModule,
		authControllers.NewAuthController,
		auth.NewAuthModule,
		articleRepositories.NewArticleRepository,
		articleServices.NewArticleService,
		articleControllers.NewArticleController,
		article.NewArticleModule,
		ProvideModules,
		New,
	)

	return nil, nil
}

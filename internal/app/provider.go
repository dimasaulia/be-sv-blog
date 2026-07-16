package app

import (
	"github.com/sv-blog/internal/modules/article"
	"github.com/sv-blog/internal/modules/auth"
	"github.com/sv-blog/internal/modules/health"
	"github.com/sv-blog/internal/modules/releasenotes"
	"github.com/sv-blog/internal/platform/config"
	"github.com/sv-blog/internal/platform/logger"
)

var _ Module = (*health.HealthModuleImpl)(nil)
var _ Module = (*releasenotes.ReleaseNoteModuleImpl)(nil)
var _ Module = (*auth.AuthModuleImpl)(nil)
var _ Module = (*article.ArticleModuleImpl)(nil)

func ProvideLogger(cfg config.Config) (*logger.Logger, error) {
	return logger.New(logger.Config{
		Level:  cfg.Logger.Level,
		LogDir: cfg.Logger.LogDir,
	})
}

func ProvideModules(healthModule *health.HealthModuleImpl, releaseNoteModule *releasenotes.ReleaseNoteModuleImpl, authModule *auth.AuthModuleImpl, articleModule *article.ArticleModuleImpl) []Module {
	return []Module{healthModule, releaseNoteModule, authModule, articleModule}
}

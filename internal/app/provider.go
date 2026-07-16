package app

import (
	"github.com/sv-blog/internal/modules/health"
	"github.com/sv-blog/internal/modules/releasenotes"
	"github.com/sv-blog/internal/platform/config"
	"github.com/sv-blog/internal/platform/logger"
)

var _ Module = (*health.HealthModuleImpl)(nil)
var _ Module = (*releasenotes.ReleaseNoteModuleImpl)(nil)

func ProvideLogger(cfg config.Config) (*logger.Logger, error) {
	return logger.New(logger.Config{
		Level:  cfg.Logger.Level,
		LogDir: cfg.Logger.LogDir,
	})
}

func ProvideModules(healthModule *health.HealthModuleImpl, releaseNoteModule *releasenotes.ReleaseNoteModuleImpl) []Module {
	return []Module{healthModule, releaseNoteModule}
}

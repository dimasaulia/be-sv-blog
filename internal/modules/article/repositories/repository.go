package repositories

import (
	"context"

	"github.com/sv-blog/internal/entities"
)

type ArticleRepository interface {
	Feed(ctx context.Context, limit uint64, offset uint64) ([]entities.Article, error)
	Find(ctx context.Context, userEmail string, status string, limit uint64, offset uint64) ([]entities.Article, error)
	FindByID(ctx context.Context, id int64) (*entities.Article, error)
	Create(ctx context.Context, article entities.Article) (*entities.Article, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Article, error)
}

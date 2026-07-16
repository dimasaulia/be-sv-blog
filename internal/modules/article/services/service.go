package services

import (
	"context"

	"github.com/sv-blog/internal/entities"
	"github.com/sv-blog/internal/modules/article/dto"
)

type ArticleService interface {
	Feed(ctx context.Context, limit uint64, offset uint64) ([]entities.Article, error)
	Find(ctx context.Context, userEmail string, status string, limit uint64, offset uint64) ([]entities.Article, error)
	FindByID(ctx context.Context, id int64, userEmail *string) (*entities.Article, error)
	Create(ctx context.Context, request dto.CreateArticleRequest, createdBy string) (*entities.Article, error)
	Update(ctx context.Context, id int64, request dto.UpdateArticleRequest, updatedBy string) (*entities.Article, error)
	Delete(ctx context.Context, id int64, updatedBy string) error
}

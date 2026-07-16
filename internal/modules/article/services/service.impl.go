package services

import (
	"context"
	"errors"
	"strings"

	"github.com/sv-blog/internal/entities"
	"github.com/sv-blog/internal/modules/article/dto"
	"github.com/sv-blog/internal/modules/article/repositories"
	"github.com/sv-blog/internal/platform/logger"
)

var ErrForbidden = errors.New("article: only the owner can modify this article")

type ArticleServiceImpl struct {
	ArticleRepository repositories.ArticleRepository
	log               *logger.LayerLogger
}

func NewArticleService(articleRepository repositories.ArticleRepository, appLogger *logger.Logger) ArticleService {
	return &ArticleServiceImpl{
		ArticleRepository: articleRepository,
		log:               appLogger.Layer("service.article"),
	}
}

func (s *ArticleServiceImpl) Feed(ctx context.Context, limit uint64, offset uint64) ([]entities.Article, error) {
	end := s.log.Start(ctx, "Feed")
	articles, err := s.ArticleRepository.Feed(ctx, limit, offset)
	filteredArticle := make([]entities.Article, 0)
	for _, article := range articles {
		filteredArticle = append(filteredArticle, entities.Article{
			ID:          article.ID,
			Title:       article.Title,
			Content:     (article.Content[50:]) + " ...",
			Category:    article.Category,
			Status:      article.Status,
			CreatedDate: article.CreatedDate,
			UpdatedDate: article.UpdatedDate,
			CreatedBy:   article.CreatedBy,
		})
	}
	end(err, "count", len(filteredArticle))
	return filteredArticle, err
}

func (s *ArticleServiceImpl) Find(ctx context.Context, userEmail string, status string, limit uint64, offset uint64) ([]entities.Article, error) {
	end := s.log.Start(ctx, "Find")

	if status != "" {
		status = strings.ToUpper(string(status[0])) + strings.ToLower(status[1:])
		errs := newValidationError()
		validStatus("status", status, errs)
		if errs.hasErrors() {
			end(errs)
			return nil, errs
		}
	}

	articles, err := s.ArticleRepository.Find(ctx, userEmail, status, limit, offset)
	filteredArticle := make([]entities.Article, 0)
	for _, article := range articles {
		filteredArticle = append(filteredArticle, entities.Article{
			ID:          article.ID,
			Title:       article.Title,
			Content:     (article.Content[200:]) + " ...",
			Category:    article.Category,
			Status:      article.Status,
			CreatedDate: article.CreatedDate,
			UpdatedDate: article.UpdatedDate,
			CreatedBy:   article.CreatedBy,
		})
	}
	end(err, "count", len(filteredArticle))
	return filteredArticle, err
}

func (s *ArticleServiceImpl) FindByID(ctx context.Context, id int64, userEmail *string) (*entities.Article, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	article, err := s.ArticleRepository.FindByID(ctx, id)
	if err != nil {
		end(err)
		return nil, err
	}

	// Kalau dia bukan artikel publish, dan email kosong atau tidak sesuai, maka user tidak berhak view data
	isOwner := userEmail != nil && article.CreatedBy != nil && *article.CreatedBy == *userEmail
	if article.Status != entities.ArticleStatusPublish && !isOwner {
		end(ErrForbidden)
		return nil, ErrForbidden
	}

	end(err)
	return article, err
}

func (s *ArticleServiceImpl) Create(ctx context.Context, request dto.CreateArticleRequest, createdBy string) (*entities.Article, error) {
	end := s.log.Start(ctx, "Create")

	errs := newValidationError()
	validTitle("title", request.Title, errs)
	validContent("content", request.Content, errs)
	validCategory("category", request.Category, errs)

	status := request.Status
	if status == "" {
		status = entities.ArticleStatusDraft
	} else {
		transformStatus := strings.ToUpper(string(status[0])) + strings.ToLower(status[1:])
		validStatus("status", transformStatus, errs)
		status = transformStatus
	}

	if errs.hasErrors() {
		end(errs)
		return nil, errs
	}

	article, err := s.ArticleRepository.Create(ctx, entities.Article{
		Title:     request.Title,
		Content:   request.Content,
		Category:  request.Category,
		Status:    status,
		CreatedBy: &createdBy,
	})
	end(err)
	return article, err
}

func (s *ArticleServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateArticleRequest, updatedBy string) (*entities.Article, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	existing, err := s.ArticleRepository.FindByID(ctx, id)
	if err != nil {
		end(err)
		return nil, err
	}

	if existing.CreatedBy == nil || *existing.CreatedBy != updatedBy {
		end(ErrForbidden)
		return nil, ErrForbidden
	}

	errs := newValidationError()
	data := map[string]any{}

	if request.Title != nil {
		validTitle("title", *request.Title, errs)
		data["title"] = *request.Title
	}
	if request.Content != nil {
		validContent("content", *request.Content, errs)
		data["content"] = *request.Content
	}
	if request.Category != nil {
		validCategory("category", *request.Category, errs)
		data["category"] = *request.Category
	}
	if request.Status != nil {
		transformStatus := strings.ToUpper(string((*request.Status)[0])) + strings.ToLower((*request.Status)[1:])
		validStatus("status", transformStatus, errs)
		data["status"] = transformStatus
	}

	if errs.hasErrors() {
		end(errs)
		return nil, errs
	}

	data["updated_by"] = updatedBy

	article, err := s.ArticleRepository.Update(ctx, id, data)
	end(err)
	return article, err
}

func (s *ArticleServiceImpl) Delete(ctx context.Context, id int64, updatedBy string) error {
	end := s.log.Start(ctx, "Delete", "id", id)

	existing, err := s.ArticleRepository.FindByID(ctx, id)
	if err != nil {
		end(err)
		return err
	}

	if existing.CreatedBy == nil || *existing.CreatedBy != updatedBy {
		end(ErrForbidden)
		return ErrForbidden
	}

	_, err = s.ArticleRepository.Update(ctx, id, map[string]any{
		"status":     entities.ArticleStatusTrash,
		"updated_by": updatedBy,
	})
	end(err)
	return err
}

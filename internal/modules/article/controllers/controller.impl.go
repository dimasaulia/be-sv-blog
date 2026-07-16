package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/sv-blog/internal/entities"
	"github.com/sv-blog/internal/modules/article/dto"
	"github.com/sv-blog/internal/modules/article/repositories"
	"github.com/sv-blog/internal/modules/article/services"
	"github.com/sv-blog/internal/platform/logger"
	"github.com/sv-blog/internal/shared/requestctx"
	"github.com/sv-blog/internal/shared/response"
)

type ArticleControllerImpl struct {
	ArticleService services.ArticleService
	response       *response.Sender
	log            *logger.LayerLogger
}

func NewArticleController(articleService services.ArticleService, sender *response.Sender, appLogger *logger.Logger) ArticleController {
	return &ArticleControllerImpl{
		ArticleService: articleService,
		response:       sender,
		log:            appLogger.Layer("controller.article"),
	}
}

func (c *ArticleControllerImpl) Feed(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	limit := parseUintPath(r, "limit", 20)
	offset := parseUintPath(r, "offset", 0)
	articles, err := c.ArticleService.Feed(r.Context(), limit, offset)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(articles))
	c.response.Success(w, r, http.StatusOK, "article.find.success", dto.ListResponse[entities.Article]{Items: articles})
}

func (c *ArticleControllerImpl) Find(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Find")

	user, ok := requestctx.CurrentUser(r.Context())
	if !ok {
		end(nil)
		c.response.Error(w, r, http.StatusUnauthorized, "auth.unauthorized", nil)
		return
	}

	limit := parseUintPath(r, "limit", 20)
	offset := parseUintPath(r, "offset", 0)
	status := r.URL.Query().Get("status")
	articles, err := c.ArticleService.Find(r.Context(), user.Email, status, limit, offset)
	var validationErr *services.ValidationError
	if errors.As(err, &validationErr) {
		end(err)
		c.response.Error(w, r, http.StatusUnprocessableEntity, "article.validation_failed", validationErr.Fields)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil, "count", len(articles))
	c.response.Success(w, r, http.StatusOK, "article.find.success", dto.ListResponse[entities.Article]{Items: articles})
}

func (c *ArticleControllerImpl) FindByID(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "FindByID")

	var userEmail *string
	if user, ok := requestctx.CurrentUser(r.Context()); ok {
		userEmail = &user.Email
	}

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "article.invalid_id", nil)
		return
	}

	article, err := c.ArticleService.FindByID(r.Context(), id, userEmail)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "article.not_found", nil)
		return
	}
	if errors.Is(err, services.ErrForbidden) {
		end(err)
		c.response.Error(w, r, http.StatusForbidden, "article.forbidden", nil)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "article.find_by_id.success", article)
}

func (c *ArticleControllerImpl) Create(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Create")

	user, ok := requestctx.CurrentUser(r.Context())
	if !ok {
		end(nil)
		c.response.Error(w, r, http.StatusUnauthorized, "auth.unauthorized", nil)
		return
	}

	var request dto.CreateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "article.invalid_payload", nil)
		return
	}

	article, err := c.ArticleService.Create(r.Context(), request, user.Email)
	var validationErr *services.ValidationError
	if errors.As(err, &validationErr) {
		end(err)
		c.response.Error(w, r, http.StatusUnprocessableEntity, "article.validation_failed", validationErr.Fields)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "article.create.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusCreated, "article.create.success", article)
}

func (c *ArticleControllerImpl) Update(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Update")

	user, ok := requestctx.CurrentUser(r.Context())
	if !ok {
		end(nil)
		c.response.Error(w, r, http.StatusUnauthorized, "auth.unauthorized", nil)
		return
	}

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "article.invalid_id", nil)
		return
	}

	var request dto.UpdateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "article.invalid_payload", nil)
		return
	}

	article, err := c.ArticleService.Update(r.Context(), id, request, user.Email)
	if errors.Is(err, repositories.ErrNotFound) {
		end(err)
		c.response.Error(w, r, http.StatusNotFound, "article.not_found", nil)
		return
	}
	if errors.Is(err, services.ErrForbidden) {
		end(err)
		c.response.Error(w, r, http.StatusForbidden, "article.forbidden", nil)
		return
	}
	var validationErr *services.ValidationError
	if errors.As(err, &validationErr) {
		end(err)
		c.response.Error(w, r, http.StatusUnprocessableEntity, "article.validation_failed", validationErr.Fields)
		return
	}
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "article.update.failed", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "article.update.success", article)
}

func (c *ArticleControllerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "Delete")

	user, ok := requestctx.CurrentUser(r.Context())
	if !ok {
		end(nil)
		c.response.Error(w, r, http.StatusUnauthorized, "auth.unauthorized", nil)
		return
	}

	id, err := parseID(r)
	if err != nil {
		end(err)
		c.response.Error(w, r, http.StatusBadRequest, "article.invalid_id", nil)
		return
	}

	if err := c.ArticleService.Delete(r.Context(), id, user.Email); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			end(err)
			c.response.Error(w, r, http.StatusNotFound, "article.not_found", nil)
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			end(err)
			c.response.Error(w, r, http.StatusForbidden, "article.forbidden", nil)
			return
		}
		end(err)
		c.response.Error(w, r, http.StatusInternalServerError, "error.internal", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "article.delete.success", nil)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

func parseUintPath(r *http.Request, key string, fallback uint64) uint64 {
	value := r.PathValue(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fallback
	}

	return parsed
}

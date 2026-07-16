package repositories

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/sv-blog/internal/entities"
	"github.com/sv-blog/internal/platform/database"
	"github.com/sv-blog/internal/platform/logger"
)

const tableName = "posts"

var ErrNotFound = errors.New("article not found")

type ArticleRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
	sb  sq.StatementBuilderType
}

func NewArticleRepository(db *database.Database, appLogger *logger.Logger) ArticleRepository {
	return &ArticleRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.article"),
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

func (r *ArticleRepositoryImpl) Feed(ctx context.Context, limit uint64, offset uint64) ([]entities.Article, error) {
	query, args, err := r.sb.Select(columns()...).
		From(tableName).
		Where(sq.Eq{"status": "Publish"}).
		OrderBy("created_date DESC", "id DESC").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	articles := make([]entities.Article, 0)
	if err := r.db.DB.SelectContext(ctx, &articles, query, args...); err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepositoryImpl) Find(ctx context.Context, userEmail string, status string, limit uint64, offset uint64) ([]entities.Article, error) {
	builder := r.sb.Select(columns()...).
		From(tableName).
		Where(sq.Eq{"created_by": userEmail})

	if status != "" {
		builder = builder.Where(sq.Eq{"status": status})
	}

	query, args, err := builder.
		OrderBy("created_date DESC", "id DESC").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	articles := make([]entities.Article, 0)
	if err := r.db.DB.SelectContext(ctx, &articles, query, args...); err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.Article, error) {
	query, args, err := r.sb.Select(columns()...).
		From(tableName).
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var article entities.Article
	if err := r.db.DB.GetContext(ctx, &article, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &article, nil
}

func (r *ArticleRepositoryImpl) Create(ctx context.Context, article entities.Article) (*entities.Article, error) {
	values := map[string]any{
		"title":      article.Title,
		"content":    article.Content,
		"category":   article.Category,
		"status":     article.Status,
		"created_by": article.CreatedBy,
	}

	query, args, err := r.sb.Insert(tableName).
		SetMap(values).
		ToSql()
	if err != nil {
		return nil, err
	}

	result, err := r.db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, id)
}

func (r *ArticleRepositoryImpl) Update(ctx context.Context, id int64, data map[string]any) (*entities.Article, error) {
	query, args, err := r.sb.Update(tableName).
		SetMap(data).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	if _, err := r.db.DB.ExecContext(ctx, query, args...); err != nil {
		return nil, err
	}

	return r.FindByID(ctx, id)
}

func columns() []string {
	return []string{
		"id",
		"title",
		"content",
		"category",
		"status",
		"created_date",
		"updated_date",
		"created_by",
		"updated_by",
	}
}

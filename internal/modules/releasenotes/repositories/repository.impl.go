package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/sv-blog/internal/entities"
	"github.com/sv-blog/internal/platform/database"
	"github.com/sv-blog/internal/platform/logger"
)

const tableName = "release_notes"

var ErrNotFound = errors.New("release note not found")

type ReleaseNoteRepositoryImpl struct {
	db  *database.Database
	log *logger.LayerLogger
	sb  sq.StatementBuilderType
}

func NewReleaseNoteRepository(db *database.Database, appLogger *logger.Logger) ReleaseNoteRepository {
	return &ReleaseNoteRepositoryImpl{
		db:  db,
		log: appLogger.Layer("repository.release_notes"),
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

func (r *ReleaseNoteRepositoryImpl) Find(ctx context.Context, limit uint64, offset uint64) ([]entities.ReleaseNote, error) {
	query, args, err := r.sb.Select(columns()...).
		From(tableName).
		Where(sq.Eq{"deleted_at": nil}).
		OrderBy("release_date DESC", "id DESC").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	releaseNotes := make([]entities.ReleaseNote, 0)
	if err := r.db.DB.SelectContext(ctx, &releaseNotes, query, args...); err != nil {
		return nil, err
	}

	return releaseNotes, nil
}

func (r *ReleaseNoteRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.ReleaseNote, error) {
	query, args, err := r.sb.Select(columns()...).
		From(tableName).
		Where(sq.Eq{"id": id, "deleted_at": nil}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var releaseNote entities.ReleaseNote
	if err := r.db.DB.GetContext(ctx, &releaseNote, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &releaseNote, nil
}

func (r *ReleaseNoteRepositoryImpl) Create(ctx context.Context, releaseNote entities.ReleaseNote) (*entities.ReleaseNote, error) {
	values := map[string]any{
		"version":      releaseNote.Version,
		"category":     releaseNote.Category,
		"title":        releaseNote.Title,
		"parent_id":    releaseNote.ParentID,
		"notes":        releaseNote.Notes,
		"release_date": releaseNote.ReleaseDate,
		"visibility":   releaseNote.Visibility,
		"is_active":    releaseNote.IsActive,
		"created_by":   releaseNote.CreatedBy,
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

func (r *ReleaseNoteRepositoryImpl) Update(ctx context.Context, id int64, data map[string]any) (*entities.ReleaseNote, error) {
	data["updated_at"] = time.Now()

	query, args, err := r.sb.Update(tableName).
		SetMap(data).
		Where(sq.Eq{"id": id, "deleted_at": nil}).
		ToSql()
	if err != nil {
		return nil, err
	}

	result, err := r.db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrNotFound
	}

	return r.FindByID(ctx, id)
}

func (r *ReleaseNoteRepositoryImpl) SoftDelete(ctx context.Context, id int64, deletedBy *int64) error {
	query, args, err := r.sb.Update(tableName).
		SetMap(map[string]any{
			"deleted_at": time.Now(),
			"deleted_by": deletedBy,
			"is_active":  int16(0),
		}).
		Where(sq.Eq{"id": id, "deleted_at": nil}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}

	return nil
}

func columns() []string {
	return []string{
		"id",
		"version",
		"category",
		"title",
		"parent_id",
		"notes",
		"release_date",
		"visibility",
		"is_active",
		"created_at",
		"updated_at",
		"deleted_at",
		"created_by",
		"updated_by",
		"deleted_by",
	}
}

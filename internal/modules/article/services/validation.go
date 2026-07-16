package services

import (
	"strings"

	"github.com/sv-blog/internal/entities"
)

const (
	titleMinLength    = 20
	contentMinLength  = 200
	categoryMinLength = 3
)

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "article: validation failed"
}

func newValidationError() *ValidationError {
	return &ValidationError{Fields: map[string]string{}}
}

func (e *ValidationError) add(field string, message string) {
	e.Fields[field] = message
}

func (e *ValidationError) hasErrors() bool {
	return len(e.Fields) > 0
}

func validTitle(field string, value string, errs *ValidationError) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		errs.add(field, "title is required")
		return
	}
	if len(trimmed) < titleMinLength {
		errs.add(field, "title must be at least 20 characters")
	}
}

func validContent(field string, value string, errs *ValidationError) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		errs.add(field, "content is required")
		return
	}
	if len(trimmed) < contentMinLength {
		errs.add(field, "content must be at least 200 characters")
	}
}

func validCategory(field string, value string, errs *ValidationError) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		errs.add(field, "category is required")
		return
	}
	if len(trimmed) < categoryMinLength {
		errs.add(field, "category must be at least 3 characters")
	}
}

func validStatus(field string, value string, errs *ValidationError) {
	switch value {
	case entities.ArticleStatusDraft, entities.ArticleStatusPublish, entities.ArticleStatusTrash:
		return
	default:
		errs.add(field, "status must be one of Draft, Publish, Trash")
	}
}

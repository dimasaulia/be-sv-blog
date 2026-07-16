package entities

import "time"

const (
	ArticleStatusDraft   = "Draft"
	ArticleStatusPublish = "Publish"
	ArticleStatusTrash   = "Trash"
)

type Article struct {
	ID          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Content     string    `db:"content" json:"content"`
	Category    string    `db:"category" json:"category"`
	Status      string    `db:"status" json:"status"`
	CreatedDate time.Time `db:"created_date" json:"created_date"`
	UpdatedDate time.Time `db:"updated_date" json:"updated_date"`
	CreatedBy   *string   `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy   *string   `db:"updated_by" json:"updated_by,omitempty"`
}

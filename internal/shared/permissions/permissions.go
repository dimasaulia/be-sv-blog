package permissions

const (
	AppCodeSVBlog = "sv-blog"

	ActionDelete = "delete"
	ActionUpdate = "update"
	ActionWrite  = "write"
	ActionRead   = "read"

	ModuleArticle = "article"
)

const (
	ArticleDelete = AppCodeSVBlog + "." + ModuleArticle + "." + ActionDelete
	ArticleUpdate = AppCodeSVBlog + "." + ModuleArticle + "." + ActionUpdate
	ArticleWrite  = AppCodeSVBlog + "." + ModuleArticle + "." + ActionWrite
	ArticleRead   = AppCodeSVBlog + "." + ModuleArticle + "." + ActionRead
)

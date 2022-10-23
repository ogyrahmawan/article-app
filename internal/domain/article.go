package domain

import (
	"article-app/pkg/database/paginator"
	"context"
	"time"

	beegoContext "github.com/beego/beego/v2/server/web/context"
	"gorm.io/gorm"
)

type Article struct {
	ID        int            `gorm:"primarykey;autoIncrement:true"`
	Author    string         `gorm:"type:text;column:author"`
	Title     string         `gorm:"type:text;column:title"`
	Body      string         `gorm:"type:text;column:body"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;autoDeleteTime"`
}

func (u *User) TableName() string {
	return "users"
}

type CreateArticleStoreRequest struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (r CreateArticleStoreRequest) ToArticle() Article {
	return Article{
		Body:   r.Body,
		Title:  r.Title,
		Author: r.Author,
	}
}

type UpdateArticleRequest struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (r UpdateArticleRequest) ToArticle() Article {
	return Article{
		ID:     r.ID,
		Body:   r.Body,
		Title:  r.Title,
		Author: r.Author,
	}
}

type GetArticleResponse struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type GetArticlesFilter struct {
	OrderBy string `json:"order_by"`
	SortBy  string `json:"sort_by"`
	Search  string `json:"search"`
	Author  string `json:"author"`
}

func (r Article) ToArticleResponse() GetArticleResponse {
	return GetArticleResponse{
		ID:     r.ID,
		Author: r.Author,
		Title:  r.Title,
		Body:   r.Body,
	}
}

type ArticleUseCase interface {
	CreateArticle(beegoCtx *beegoContext.Context, data CreateArticleStoreRequest) (*GetArticleResponse, error)
	GetArticles(beegoCtx *beegoContext.Context, page, limit, offset int, filter GetArticlesFilter) (result *paginator.Paginator, err error)
	GetArticleById(beegoCtx *beegoContext.Context, id int) (*GetArticleResponse, error)
	UpdateArticle(beegoCtx *beegoContext.Context, body UpdateArticleRequest, id int) (*GetArticleResponse, error)
	DeleteArticle(beegoCtx *beegoContext.Context, id int) error
}

type ArticleRepository interface {
	Store(ctx context.Context, tx *gorm.DB, data Article) (int, error)
	FetchWithFilterAndPagination(ctx context.Context, page, limit int, offset int, order string, fields, associate, filter []string, model interface{}, args ...interface{}) (*paginator.Paginator, error)
	FindByID(ctx context.Context, id int) (*Article, error)
	Update(ctx context.Context, body Article, id int) error
	Delete(ctx context.Context, id int) error
	DB() *gorm.DB
}

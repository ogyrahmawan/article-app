package repository

import (
	"article-app/internal/domain"
	"article-app/pkg/database/paginator"
	"context"
	"strings"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) domain.ArticleRepository {
	return &ArticleRepository{
		db: db,
	}
}

func (ar ArticleRepository) DB() *gorm.DB {
	return ar.db
}

func (ar ArticleRepository) Store(ctx context.Context, tx *gorm.DB, data domain.Article) (int, error) {
	err := tx.WithContext(ctx).Create(&data).Error
	if err != nil {
		return 0, err
	}
	return data.ID, nil
}

func (ar ArticleRepository) FetchWithFilterAndPagination(ctx context.Context, page, limit int, offset int, order string, fields, associate, filter []string, model interface{}, args ...interface{}) (*paginator.Paginator, error) {
	p := paginator.NewPaginator(ar.db, page, limit, model)
	if err := p.FindWithFilter(ctx, order, fields, associate, filter, args...).Select(strings.Join(fields, ",")).Error; err != nil {
		return p, err
	}
	return p, nil
}

func (ar ArticleRepository) FindByID(ctx context.Context, id int) (*domain.Article, error) {
	var entity domain.Article
	err := ar.db.WithContext(ctx).First(&entity, "id =?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (ar ArticleRepository) Update(ctx context.Context, data domain.Article, id int) error {
	err := ar.db.WithContext(ctx).Where("articles.id = ?", id).Updates(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (ar ArticleRepository) Delete(ctx context.Context, id int) error {
	err := ar.db.WithContext(ctx).Exec("delete from articles where id =?", id).Error
	if err != nil {
		return err
	}
	return nil
}

package usecase

import (
	"article-app/internal/domain"
	"article-app/pkg/database/paginator"
	"article-app/pkg/jwt"
	"context"
	"fmt"
	"time"

	beegoContext "github.com/beego/beego/v2/server/web/context"
	"gorm.io/gorm"
)

type articleUseCase struct {
	contextTimeout time.Duration
	articleRepo    domain.ArticleRepository
	jwtAuth        jwt.JWT
	expireToken    int
}

func NewArticleUseCase(timeout time.Duration, ur domain.ArticleRepository, jwtAuth jwt.JWT, expireToken int) domain.ArticleUseCase {
	return &articleUseCase{
		contextTimeout: timeout,
		articleRepo:    ur,
		jwtAuth:        jwtAuth,
		expireToken:    expireToken,
	}
}

func (auc articleUseCase) CreateArticle(beegoCtx *beegoContext.Context, body domain.CreateArticleStoreRequest) (*domain.GetArticleResponse, error) {
	var id int
	ctx, cancel := context.WithTimeout(beegoCtx.Request.Context(), auc.contextTimeout)
	defer cancel()

	auc.articleRepo.DB().Transaction(func(tx *gorm.DB) error {
		if articleId, err := auc.articleRepo.Store(ctx, tx, body.ToArticle()); err != nil {
			return err
		} else {
			//set returning id from db
			id = articleId
		}
		return nil
	})
	data, err := auc.articleRepo.FindByID(ctx, id)
	res := data.ToArticleResponse()
	return &res, err
}

func (auc articleUseCase) GetArticles(beegoCtx *beegoContext.Context, page, limit, offset int, filter domain.GetArticlesFilter) (result *paginator.Paginator, err error) {
	ctx, cancel := context.WithTimeout(beegoCtx.Request.Context(), auc.contextTimeout)
	defer cancel()
	var entities []domain.Article

	var where []string

	if filter.Author != "" {
		where = append(where, fmt.Sprintf("articles.author = '%s'", filter.Author))
	}

	if filter.Search != "" {
		where = append(
			where,
			fmt.Sprintf(`articles.title like "%%%s%%" OR articles.body like "%%%s%%"`, filter.Search, filter.Search),
		)
	}

	paging, err := auc.articleRepo.FetchWithFilterAndPagination(ctx,
		page,
		limit,
		offset,
		fmt.Sprintf("articles.%s %s", filter.SortBy, filter.OrderBy),
		[]string{
			"articles.id",
			"articles.title",
			"articles.author",
			"articles.body",
		}, nil, where, &entities, nil,
	)
	if err != nil {
		return nil, err
	}

	var dataList = make([]domain.GetArticleResponse, len(entities))

	for k, v := range entities {
		dataList[k] = v.ToArticleResponse()
	}
	paging.Records = dataList
	return paging, nil
}

func (auc articleUseCase) GetArticleById(beegoCtx *beegoContext.Context, id int) (*domain.GetArticleResponse, error) {
	ctx, cancel := context.WithTimeout(beegoCtx.Request.Context(), auc.contextTimeout)
	defer cancel()

	data, err := auc.articleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := data.ToArticleResponse()
	return &res, err
}

func (auc articleUseCase) UpdateArticle(beegoCtx *beegoContext.Context, body domain.UpdateArticleRequest, id int) (*domain.GetArticleResponse, error) {
	ctx, cancel := context.WithTimeout(beegoCtx.Request.Context(), auc.contextTimeout)
	defer cancel()

	data := body.ToArticle()
	err := auc.articleRepo.Update(ctx, data, id)
	if err != nil {
		return nil, err
	}

	article, err := auc.articleRepo.FindByID(ctx, id)
	res := article.ToArticleResponse()
	return &res, err
}

func (auc articleUseCase) DeleteArticle(beegoCtx *beegoContext.Context, id int) error {
	ctx, cancel := context.WithTimeout(beegoCtx.Request.Context(), auc.contextTimeout)
	defer cancel()

	found, err := auc.articleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if found != nil {
		err = auc.articleRepo.Delete(ctx, id)
		if err != nil {
			return err
		}
	}
	return nil
}

package http

import (
	"article-app/internal"
	"article-app/internal/domain"
	"article-app/pkg/database/paginator"
	"article-app/pkg/jwt"
	"article-app/pkg/response"
	"net/http"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

type articleHandler struct {
	internal.BaseController
	response.ApiResponse
	ArticleUseCase domain.ArticleUseCase
	JwtAuth        jwt.JWT
}

func NewArticleHandler(useCase domain.ArticleUseCase, jwt jwt.JWT) {
	pHandler := &articleHandler{
		ArticleUseCase: useCase,
		JwtAuth:        jwt,
	}
	beego.Router("/api/v1/cms/article", pHandler, "post:CreateArticle")
	beego.Router("/api/v1/cms/article", pHandler, "get:GetArticles")
	beego.Router("/api/v1/cms/article/:id", pHandler, "get:GetArticleById")
	beego.Router("/api/v1/cms/article/:id", pHandler, "patch:UpdateArticle")
	beego.Router("/api/v1/cms/article/:id", pHandler, "delete:DeleteArticle")
}

func (h *articleHandler) Prepare() {
	// check user access when needed
	h.SetLangVersion()
}

func (h *articleHandler) CreateArticle() {
	var request domain.CreateArticleStoreRequest
	if err := h.BindJSON(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationCodeError, domain.ErrorCodeText(domain.ApiValidationCodeError, h.Locale.Lang), err)
		return
	}

	data, err := h.ArticleUseCase.CreateArticle(h.Ctx, request)
	if err != nil {
		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), err)
		return
	}
	h.Ok(h.Ctx, h.Tr("message.success"), data)
	return
}

func (h *articleHandler) GetArticles() {
	pageSize, page, err := domain.PaginationQueryParamValidation(h.Ctx.Input.Query("pageSize"), h.Ctx.Input.Query("page"))
	if err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.QueryParamInvalidCode, domain.ErrorCodeText(domain.QueryParamInvalidCode, h.Locale.Lang), err)
		return
	}

	filter := domain.GetArticlesFilter{
		OrderBy: h.Ctx.Input.Query("order_by"),
		SortBy:  h.Ctx.Input.Query("sort_by"),
		Search:  h.Ctx.Input.Query("search"),
		Author:  h.Ctx.Input.Query("author"),
	}

	limit, page, offset := paginator.Pagination(page, pageSize)

	result, err := h.ArticleUseCase.GetArticles(h.Ctx, page, limit, offset, filter)

	if err != nil {
		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), err)
		return
	}
	h.Ok(h.Ctx, h.Tr("message.success"), result)

	return
}

func (h *articleHandler) GetArticleById() {
	pathParam, err := strconv.Atoi(h.Ctx.Input.Param(":id"))
	if err != nil || pathParam < 1 {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.PathParamInvalidCode, domain.ErrorCodeText(domain.PathParamInvalidCode, h.Locale.Lang), err)
		return
	}

	result, err := h.ArticleUseCase.GetArticleById(h.Ctx, pathParam)
	if err != nil {
		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), err)
		return
	}
	h.Ok(h.Ctx, h.Tr("message.success"), result)

	return
}

func (h *articleHandler) UpdateArticle() {
	pathParam, err := strconv.Atoi(h.Ctx.Input.Param(":id"))
	if err != nil || pathParam < 1 {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.PathParamInvalidCode, domain.ErrorCodeText(domain.PathParamInvalidCode, h.Locale.Lang), err)
		return
	}

	var request domain.UpdateArticleRequest
	h.BindJSON(&request)

	if err := h.BindJSON(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationCodeError, domain.ErrorCodeText(domain.ApiValidationCodeError, h.Locale.Lang), err)
		return
	}

	data, err := h.ArticleUseCase.UpdateArticle(h.Ctx, request, pathParam)
	if err != nil {
		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), err)
		return
	}
	h.Ok(h.Ctx, h.Tr("message.success"), data)
	return
}

func (h *articleHandler) DeleteArticle() {
	pathParam, err := strconv.Atoi(h.Ctx.Input.Param(":id"))
	if err != nil || pathParam < 1 {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.PathParamInvalidCode, domain.ErrorCodeText(domain.PathParamInvalidCode, h.Locale.Lang), err)
		return
	}

	h.ArticleUseCase.DeleteArticle(h.Ctx, pathParam)
	if err != nil {
		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), err)
		return
	}
	h.Ok(h.Ctx, h.Tr("message.success"), nil)
	return
}

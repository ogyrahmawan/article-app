package middlewares

import (
	"article-app/internal/domain"
	"article-app/pkg/helper"
	"article-app/pkg/jwt"
	"article-app/pkg/response"
	"net/http"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

type JwtConfig struct {
	Skipper Skipper
	response.ApiResponse
}

func NewJwtMiddleware() *JwtConfig {
	return &JwtConfig{Skipper: func(ctx *context.Context) bool {
		if strings.EqualFold(ctx.Request.URL.Path, "/api/v1/cms/auth/token") {
			return true
		}
		if strings.EqualFold(ctx.Request.URL.Path, "/api/v1/cms/auth/register") {
			return true
		}
		if strings.EqualFold(ctx.Request.URL.Path, "/api/v1/cms/user/login") {
			return true
		}
		if strings.EqualFold(ctx.Request.URL.Path, "/api/v1/employee/auth/login") {
			return true
		}
		return false
	}}
}

func (r *JwtConfig) JwtMiddleware(jwtAuth jwt.JWT) beego.FilterChain {
	return func(next beego.FilterFunc) beego.FilterFunc {
		return func(ctx *context.Context) {
			if r.Skipper(ctx) {
				next(ctx)
				return
			}
			if ctx.Request.Method == "OPTIONS" {
				next(ctx)
				return
			}
			if middlewareRequest, err := jwtAuth.Middleware(ctx.Request); err != nil {
				switch {
				case jwt.IsInvalidToken(err):
					r.ResponseError(ctx, http.StatusUnauthorized, domain.InvalidTokenCodeError, domain.ErrorCodeText(domain.InvalidTokenCodeError, helper.GetLangVersion(ctx)), err)
					return
				case jwt.IsExpiredToken(err):
					r.ResponseError(ctx, http.StatusUnauthorized, domain.ExpiredTokenCodeError, domain.ErrorCodeText(domain.ExpiredTokenCodeError, helper.GetLangVersion(ctx)), err)
					return
				case jwt.IsMissingToken(err):
					r.ResponseError(ctx, http.StatusUnauthorized, domain.MissingTokenCodeError, domain.ErrorCodeText(domain.MissingTokenCodeError, helper.GetLangVersion(ctx)), err)
					return
				case jwt.IsAuthElsewhere(err):
					r.ResponseError(ctx, http.StatusUnauthorized, domain.AuthElseWhereCodeError, domain.ErrorCodeText(domain.AuthElseWhereCodeError, helper.GetLangVersion(ctx)), err)
					return
				default:
					r.ResponseError(ctx, http.StatusUnauthorized, domain.UnauthorizedCodeError, domain.ErrorCodeText(domain.UnauthorizedCodeError, helper.GetLangVersion(ctx)), err)
					return
				}
			} else {
				ctx.Request = middlewareRequest
				next(ctx)
			}
		}
	}
}

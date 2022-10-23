package http

import (
	"article-app/internal"
	"article-app/internal/domain"
	"article-app/pkg/jwt"
	"article-app/pkg/response"
	"errors"
	"fmt"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type UserHandler struct {
	internal.BaseController
	response.ApiResponse
	UserUseCase domain.UserUseCase
	JwtAuth     jwt.JWT
}

func NewUserHandler(useCase domain.UserUseCase, jwt jwt.JWT) {
	pHandler := &UserHandler{
		UserUseCase: useCase,
		JwtAuth:     jwt,
	}
	beego.Router("/api/v1/cms/user/login", pHandler, "post:RequestToken")
}

func (h *UserHandler) Prepare() {
	// check user access when needed
	h.SetLangVersion()
}

// RequestToken
// @Title RequestToken
// @Summary Generate JWT Token
// @Produce json
// @Tags User Auth
// @Success 200 {object} swagger.BaseResponse
// @Failure 408 {object} swagger.RequestTimeoutResponse{errors=[]object,data=object}
// @Failure 500 {object} swagger.InternalServerErrorResponse{errors=[]object,data=object}
// @Param Accept-Language header string false "lang"
// @Router /v1/cms/user/login [post]
func (h *UserHandler) RequestToken() {
	username, password, ok := h.Ctx.Request.BasicAuth()
	fmt.Println(username, password)
	if !ok {
		fmt.Println(`something went wrong`)
	}

	result, err := h.UserUseCase.Login(h.Ctx, username, password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidEmailPassword) {
			h.ResponseErrorWithData(h.Ctx, http.StatusBadRequest, domain.InvalidEmailPassword, domain.ErrorCodeText(domain.InvalidEmailPassword, h.Locale.Lang), err, result)
			return
		}

		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), err)
		return
	}
	h.Ok(h.Ctx, h.Tr("message.success"), result)
	return
}

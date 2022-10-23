package usecase

import (
	"article-app/internal/domain"
	"context"
	"time"

	"article-app/pkg/jwt"

	beegoContext "github.com/beego/beego/v2/server/web/context"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	contextTimeout time.Duration
	userRepository domain.UserRepository
	jwtAuth        jwt.JWT
	expireToken    int
}

func NewUserUseCase(timeout time.Duration, ur domain.UserRepository, jwtAuth jwt.JWT, expireToken int) domain.UserUseCase {
	return &userUseCase{
		contextTimeout: timeout,
		userRepository: ur,
		jwtAuth:        jwtAuth,
		expireToken:    expireToken,
	}
}

func (usc userUseCase) Login(beegoCtx *beegoContext.Context, email, password string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(beegoCtx.Request.Context(), usc.contextTimeout)
	defer cancel()

	result, err := usc.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err != nil {
		return nil, domain.ErrInvalidEmailPassword
	}

	token, err := usc.jwtAuth.Ctx(ctx).GenerateToken(jwt.Payload{"uid": result.Id, "email": result.Email}, beegoCtx.Request.Host, usc.expireToken)
	if err != nil {
		return nil, err
	}

	res := new(domain.UserLoginResponse)
	res.Token = token.Token
	res.ExpiredAt = token.ExpiredAt.String()
	res.User = domain.UserLogin{
		Id:    int(result.Id),
		Email: result.Email,
	}

	return res, nil
}

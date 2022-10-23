package domain

import (
	"errors"
	"fmt"

	"github.com/beego/i18n"
)

const (
	InvalidEmailPassword      = "ART-00001"
	InvalidTokenCodeError     = "ART-00002"
	ExpiredTokenCodeError     = "ART-00003"
	MissingTokenCodeError     = "ART-00004"
	AuthElseWhereCodeError    = "ART-00005"
	UnauthorizedCodeError     = "ART-00006"
	ResourceNotFoundCodeError = "ART-00007"
	ServerErrorCode           = "ART-00008"
	ApiValidationCodeError    = "ART-00009"
	RequestTimeoutCodeError   = "ART-00010"

	//Url Query & Param error
	QueryParamInvalidCode = "ART-API-001"
	PathParamInvalidCode  = "ART-API-002"
)

var (
	//query param invalid
	ErrQueryParamInvalid = errors.New("query param is invalid")

	//login auth validation
	ErrInvalidEmailPassword = errors.New("email Tidak Terdaftar atau kata sandi anda salah")
)

func ErrorCodeText(code, locale string, args ...interface{}) string {
	fmt.Println(locale, "isi locale")
	switch code {
	case InvalidEmailPassword:
		return i18n.Tr(locale, "message.errorInvalidEmailPassword", args)
	case QueryParamInvalidCode:
		return i18n.Tr(locale, "message.errorQueryParamInvalid", args)
	case PathParamInvalidCode:
		return i18n.Tr(locale, "message.errorPathParamInvalid", args)
	case ApiValidationCodeError:
		return i18n.Tr(locale, "message.errorValidation", args)
	case ResourceNotFoundCodeError:
		return i18n.Tr(locale, "message.errorResourceNotFound", args)
	case ServerErrorCode:
		return i18n.Tr(locale, "message.errorServerError", args)
	case RequestTimeoutCodeError:
		return i18n.Tr(locale, "message.errorRequestTimeout", args)
	case MissingTokenCodeError:
		return i18n.Tr(locale, "message.errorMissingToken", args)
	default:
		return ""
	}
}

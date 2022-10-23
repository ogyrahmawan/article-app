package response

import (
	"fmt"
	"net/http"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

type ApiResponse struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Errors    []Errors    `json:"errors"`
	RequestId string      `json:"request_id"`
	TimeStamp string      `json:"time_stamp"`
}

type Errors struct {
	Field       string `json:"field"`
	Description string `json:"message"`
}

func (r ApiResponse) ResponseError(ctx *context.Context, httpStatus int, errorCode string, message string, err error) error {
	var apiResponse ApiResponse
	var errorValidations []Errors = nil

	ctx.Output.SetStatus(httpStatus)

	apiResponse.RequestId = ctx.ResponseWriter.ResponseWriter.Header().Get("X-REQUEST-ID")
	apiResponse.Code = errorCode
	apiResponse.Message = message
	apiResponse.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
	apiResponse.Errors = errorValidations

	return ctx.Output.JSON(apiResponse, beego.BConfig.RunMode != "prod", false)
}

func (r ApiResponse) ResponseErrorWithData(ctx *context.Context, httpStatus int, errorCode string, message string, err error, data interface{}) error {
	fmt.Println(message)
	var apiResponse ApiResponse
	var errorValidations []Errors = nil

	ctx.Output.SetStatus(httpStatus)

	apiResponse.RequestId = ctx.Input.Header("X-REQUEST-ID")
	apiResponse.Code = errorCode
	apiResponse.Message = message
	apiResponse.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
	apiResponse.Errors = errorValidations
	apiResponse.Data = data

	return ctx.Output.JSON(apiResponse, beego.BConfig.RunMode != "prod", false)
}

func (r ApiResponse) Ok(ctx *context.Context, message string, data interface{}) error {
	ctx.Output.SetStatus(http.StatusOK)

	return ctx.Output.JSON(ApiResponse{
		Code:      http.StatusText(http.StatusOK),
		RequestId: ctx.ResponseWriter.ResponseWriter.Header().Get("X-REQUEST-ID"),
		Message:   message,
		Data:      data,
		TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
	}, beego.BConfig.RunMode != "prod", false)
}

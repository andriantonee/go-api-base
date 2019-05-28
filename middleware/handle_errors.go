package middleware

import (
	"errors"
	"go-api-base/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrInternalServer       = errors.New("Internal Server Error")
	ErrUnauthorized         = errors.New("Unauthorized")
	ErrBadRequest           = errors.New("Bad Request")
	ErrBadRequestWithDetail = errors.New("Bad Request")
)

func HandleErrors(ginContext *gin.Context) {
	ginContext.Next()

	err := ginContext.Errors.ByType(gin.ErrorTypePublic).Last()
	if err != nil {
		var httpError response.HttpError

		switch err.Err {
		case ErrUnauthorized:
			httpError = response.NewHttpError(
				http.StatusUnauthorized,
				err.Error(),
			)
		case ErrBadRequestWithDetail:
			if detail, ok := err.Meta.([]string); ok {
				httpError = response.NewHttpErrorWithDetail(
					http.StatusBadRequest,
					err.Error(),
					detail,
				)
			} else {
				panic("Meta (type []string) is required, middleware.HandleErrors case ErrBadRequestWithDetail")
			}
		case ErrBadRequest:
			if message, ok := err.Meta.(string); ok {
				httpError = response.NewHttpError(
					http.StatusBadRequest,
					message,
				)
			} else {
				panic("Meta (type string) is required, middleware.HandleErrors case ErrBadRequest")
			}
		default:
			httpError = response.NewHttpError(
				http.StatusInternalServerError,
				"Internal Server Error",
			)
		}

		ginContext.JSON(200, httpError)
	}
}

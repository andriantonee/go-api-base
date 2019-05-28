package router

import (
	"fmt"
	"go-api-base/app/auth"
	"go-api-base/app/user"
	"go-api-base/middleware"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
)

func New(userService user.Service, authService auth.Service) *gin.Engine {
	// f, _ := os.Create("gin.log")
	// gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()

	// router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	// 	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
	// 		param.ClientIP,
	// 		param.TimeStamp.Format(time.RFC1123),
	// 		param.Method,
	// 		param.Path,
	// 		param.Request.Proto,
	// 		param.StatusCode,
	// 		param.Latency,
	// 		param.Request.UserAgent(),
	// 	)
	// }))
	// router.Use(gin.Recovery())

	router.Use(middleware.HandleErrors)

	userHandlerInstance := userHandler{
		userService: userService,
		authService: authService,
	}

	api := router.Group("/api")
	{
		api.GET("/user", userHandlerInstance.index)
		api.POST("/user", userHandlerInstance.register)
		api.POST("/auth/login", userHandlerInstance.login)
	}

	return router
}

func internalServerError(ginContext *gin.Context, err error) {
	ginContext.Error(err)

	ginErr := ginContext.Error(middleware.ErrInternalServer)
	ginErr.SetType(gin.ErrorTypePublic)
}

func unauthorized(ginContext *gin.Context, err error) {
	ginContext.Error(err)

	ginErr := ginContext.Error(middleware.ErrUnauthorized)
	ginErr.SetType(gin.ErrorTypePublic)
}

func badRequestValidationMessage(
	ginContext *gin.Context,
	request interface{},
	validationErrors validator.ValidationErrors,
) {
	detail := make([]string, 0, 0)

	for _, e := range validationErrors {
		var (
			message string
			field   = strings.ToLower(e.NameNamespace)
		)

		switch e.Tag {
		case "required":
			message = fmt.Sprintf(
				"%s is required",
				field,
			)
		case "max":
			message = fmt.Sprintf(
				"%s cannot be longer than %s",
				field,
				e.Param,
			)
		case "min":
			message = fmt.Sprintf(
				"%s must be longer than %s",
				field,
				e.Param,
			)
		case "email":
			message = fmt.Sprintf(
				"%s has invalid email format",
				field,
			)
		default:
			message = fmt.Sprintf(
				"%s is not valid",
				field,
			)
		}

		detail = append(detail, message)
	}

	ginErr := ginContext.Error(middleware.ErrBadRequestWithDetail)
	ginErr.SetType(gin.ErrorTypePublic)
	ginErr.SetMeta(detail)
}

func badRequestEOF(ginContext *gin.Context) {
	ginErr := ginContext.Error(middleware.ErrBadRequestWithDetail)
	ginErr.SetType(gin.ErrorTypePublic)
	ginErr.SetMeta([]string{"body is empty"})
}

func badRequest(ginContext *gin.Context, err error) {
	ginErr := ginContext.Error(middleware.ErrBadRequest)
	ginErr.SetType(gin.ErrorTypePublic)
	ginErr.SetMeta(err.Error())
}

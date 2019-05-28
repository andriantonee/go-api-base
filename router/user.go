package router

import (
	"fmt"
	"go-api-base/app/auth"
	"go-api-base/app/user"
	"go-api-base/model"
	"io"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func (userHandler *userHandler) register(ginContext *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email,max=255"`
		Password string `json:"password" binding:"required,min=4,max=20"`
		Name     string `json:"name" binding:"required,max=255"`
	}

	if err := ginContext.ShouldBind(&request); err != nil {
		if err != io.EOF {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				badRequestValidationMessage(
					ginContext,
					request,
					validationErrors,
				)
				return
			}
		}
		badRequestEOF(ginContext)
		return
	}

	user := model.NewUser(request.Email, request.Password, request.Name)

	if err := userHandler.userService.Register(user); err != nil {
		badRequest(ginContext, err)
		return
	}

	tokenString, err := userHandler.authService.NewIdentifier(user.UserID)
	if err != nil {
		internalServerError(ginContext, err)
		return
	}
	fmt.Println(tokenString)
}

func (userHandler *userHandler) login(ginContext *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ginContext.ShouldBind(&request); err != nil {
		if err != io.EOF {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				badRequestValidationMessage(
					ginContext,
					request,
					validationErrors,
				)
				return
			}
		}
		badRequestEOF(ginContext)
		return
	}

	userID, err := userHandler.userService.FindUserID(
		request.Email,
		request.Password,
	)
	if err != nil {
		badRequest(ginContext, err)
		return
	}

	tokenString, err := userHandler.authService.NewIdentifier(*userID)
	if err != nil {
		internalServerError(ginContext, err)
		return
	}
	fmt.Println(tokenString)
}

func (userHandler *userHandler) index(ginContext *gin.Context) {
	tokenString := ginContext.GetHeader("Authorization")

	userID, err := userHandler.authService.Authorize(tokenString)
	if err != nil {
		unauthorized(ginContext, err)
		return
	}
	fmt.Println(*userID)
}

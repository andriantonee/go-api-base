package router

import (
	"fmt"
	"go-api-base/app/user"
	"go-api-base/model"

	"net/http"

	"github.com/gin-gonic/gin"
)

type userRegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=4,max=20"`
	Name     string `json:"name" binding:"required,max=255"`
}

type userLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type authTokenResponse struct {
	Token string `json:"token"`
}

type userHandler struct {
	userService user.Service
}

func (userHandler *userHandler) register(ginContext *gin.Context) {
	var request userRegisterRequest

	if err := ginContext.ShouldBind(&request); err != nil {
		// if err != io.EOF {
		// 	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// 		badRequestValidator(ginContext, request, &validationErrors)

		// 		return
		// 	}
		// }
		ginContext.JSON(
			http.StatusOK,
			newErrorResponse(http.StatusBadRequest, "Bad Request"),
		)

		return
	}

	user := model.User{
		UserID:   model.NewUserID(),
		Email:    request.Email,
		Password: model.NewPassword(request.Password),
		Name:     request.Name,
	}

	token, err := userHandler.userService.Register(user)
	if err != nil {
		fmt.Println(err)
		ginContext.JSON(
			http.StatusInternalServerError,
			newErrorResponse(
				http.StatusInternalServerError,
				"Internal Server Error",
			),
		)

		return
	}

	response := &successResponse{
		Code: http.StatusOK,
		Data: authTokenResponse{
			Token: token,
		},
	}
	ginContext.JSON(http.StatusOK, response)
}

func (userHandler *userHandler) login(ginContext *gin.Context) {
	var request userLoginRequest

	if err := ginContext.ShouldBind(&request); err != nil {
		ginContext.JSON(
			http.StatusOK,
			newErrorResponse(http.StatusBadRequest, "Bad Request"),
		)

		return
	}

	token, err := userHandler.userService.Login(
		request.Email,
		request.Password,
	)
	if err != nil {
		ginContext.JSON(
			http.StatusInternalServerError,
			newErrorResponse(
				http.StatusInternalServerError,
				"Internal Server Error",
			),
		)

		return
	}

	response := &successResponse{
		Code: http.StatusOK,
		Data: authTokenResponse{
			Token: token,
		},
	}
	ginContext.JSON(http.StatusOK, response)
}

package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (u *userHandler) Register(context *gin.Context) {
	input := user.RegisterInputUser{}
	err := context.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Register failed", http.StatusUnprocessableEntity, "error", errorMessage)
		context.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	Newuser, err := u.userService.RegisterUser(input)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Register failed", http.StatusBadRequest, "error", errorMessage)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	FormatUser := user.FormatUser(Newuser, "token")

	response := helper.APIResponse("Register succes", http.StatusOK, "succes", FormatUser)

	context.JSON(http.StatusOK, response)
}

func (u *userHandler) Login(context *gin.Context) {
	var input user.LoginInput
	err := context.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"error": errors}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		context.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedUser, err := u.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", errorMessage)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	data := user.FormatUser(loggedUser, "trdsfsdfs")
	response := helper.APIResponse("Successfully Loggein", http.StatusOK, "success", data)
	context.JSON(http.StatusOK, response)
}

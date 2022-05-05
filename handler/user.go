package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
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

	token, err := u.authService.GenerateToken(Newuser.ID)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Register failed", http.StatusBadRequest, "error", errorMessage)
		context.JSON(http.StatusBadRequest, response)
		return
	}
	FormatUser := user.FormatUser(Newuser, token)

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

	token, err := u.authService.GenerateToken(loggedUser.ID)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Register failed", http.StatusBadRequest, "error", errorMessage)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	data := user.FormatUser(loggedUser, token)
	response := helper.APIResponse("Successfully Loggein", http.StatusOK, "success", data)
	context.JSON(http.StatusOK, response)
}

func (u *userHandler) CheckEmail(context *gin.Context) {
	var input user.CheckEmailInput
	err := context.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Check email failed", http.StatusUnprocessableEntity, "error", errorMessage)
		context.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := u.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "server error"}
		response := helper.APIResponse("Check email failed", http.StatusBadRequest, "error", errorMessage)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	metaMessage := "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	data := gin.H{"is_available": isEmailAvailable}
	response := helper.APIResponse(metaMessage, http.StatusBadRequest, "error", data)
	context.JSON(http.StatusBadRequest, response)
}

func (u *userHandler) UploadAvatar(context *gin.Context) {
	file, err := context.FormFile("avatar")
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Upload avatar failed", http.StatusBadRequest, "error", errorMessage)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := context.MustGet("currentUser").(user.User)
	userId := currentUser.ID
	path := fmt.Sprintf("images/%d-%s", userId, file.Filename)
	err = context.SaveUploadedFile(file, path)
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Upload avatar failed", http.StatusBadRequest, "error", errorMessage)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = u.userService.SaveAvatar(userId, path)
	if err != nil {
		errorMessage := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Upload avatar failed", http.StatusBadRequest, "error", errorMessage)
		context.JSON(http.StatusBadRequest, response)
		return
	}
	data := gin.H{"is_uploaded": false}
	response := helper.APIResponse("Avatar successfully uploaded", http.StatusOK, "success", data)
	context.JSON(http.StatusOK, response)
}

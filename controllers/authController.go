package controllers

import (
	"github.com/gmkanat/Go-Shop/initializers"
	"github.com/gmkanat/Go-Shop/models"
	"github.com/gmkanat/Go-Shop/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

// SignUpUser [...] SignUp User
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "fail", "message": err.Error(),
		})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "fail", "message": "Passwords do not match",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status": "error", "message": err.Error(),
		})
		return
	}

	now := time.Now()
	newUser := models.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
		RoleId:    payload.RoleId,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status": "error", "message": "Something bad happened",
		})
		return
	}

	ac.DB.Save(newUser)
	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

// SignInUser [...] SignIn User
func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "fail", "message": err.Error(),
		})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "fail", "message": "Invalid email or Password",
		})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "fail", "message": "Invalid email or Password",
		})
		return
	}

	config, _ := initializers.LoadConfig(".")

	token, err := utils.GenerateToken(
		config.AccessTokenExpiresIn, user.ID, config.TokenSecret,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "fail", "message": err.Error(),
		})
		return
	}

	ctx.SetCookie("token", token, config.AccessTokenMaxAge*60,
		"/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/",
		"localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

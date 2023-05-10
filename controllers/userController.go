package controllers

import (
	"github.com/gmkanat/Go-Shop/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	userResponse := &models.UserResponse{
		ID:        currentUser.ID,
		Name:      currentUser.Name,
		Email:     currentUser.Email,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success", "data": gin.H{"user": userResponse},
	})
}

func (uc *UserController) CancelOrder(ctx *gin.Context) {
	order := ctx.MustGet("currentOrder").(models.Order)

	order.Status = "canceled"

	result := uc.DB.Save(order)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}
	NewOrderResponce := models.OrderResponse{
		ID:     order.ID,
		Status: order.Status,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "order": NewOrderResponce})
}

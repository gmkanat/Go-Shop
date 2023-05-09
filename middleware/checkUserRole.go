package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gmkanat/Go-Shop/models"
	"gorm.io/gorm"
	"net/http"
)

func CheckUserRole(DB *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := ctx.MustGet("currentUser").(models.User)
		DB.First(&currentUser.Role, currentUser.RoleId)
		if currentUser.Role.Name != "seller" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status": "fail", "message": "You have not access",
			})
			return
		}
		ctx.Next()
	}
}

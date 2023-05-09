package middleware

import (
	"github.com/gmkanat/Go-Shop/models"
	"net/http"
)

func CheckUserRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := ctx.MustGet("currentUser").(models.User)
		if currentUser.Role.name != "seller" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status": "fail", "message": "You have not access",
			})
			return
		}
		ctx.Next()
	}
}

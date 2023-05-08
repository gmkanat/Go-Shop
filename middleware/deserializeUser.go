package middleware

import (
	"fmt"
	"github.com/gmkanat/Go-Shop/initializers"
	"github.com/gmkanat/Go-Shop/models"
	"github.com/gmkanat/Go-Shop/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		cookie, err := ctx.Cookie("token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else if err == nil {
			token = cookie
		}

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "fail", "message": "You are not logged in",
			})
			return
		}

		config, _ := initializers.LoadConfig(".")
		sub, err := utils.ValidateToken(token, config.TokenSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "fail", "message": err.Error(),
			})
			return
		}

		var user models.User
		result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  "fail",
				"message": "the user belonging to this token no logger exists",
			})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}

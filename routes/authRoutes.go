package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gmkanat/Go-Shop/controllers"
	"github.com/gmkanat/Go-Shop/middleware"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	router.POST("/register", rc.authController.SignUpUser)
	router.POST("/login", rc.authController.SignInUser)
	router.GET("/logout", middleware.DeserializeUser(), rc.authController.LogoutUser)
}

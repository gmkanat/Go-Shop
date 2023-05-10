package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gmkanat/Go-Shop/controllers"
	"github.com/gmkanat/Go-Shop/middleware"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("users")
	router.GET("/me", middleware.DeserializeUser(), uc.userController.GetMe)
	router.POST("/:user_id/orders/:order_id/cancel", middleware.DeserializeUser(), uc.userController.CancelOrder)
}

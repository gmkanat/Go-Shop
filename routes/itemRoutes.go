package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gmkanat/Go-Shop/controllers"
	"github.com/gmkanat/Go-Shop/middleware"
)

type ItemRouteController struct {
	itemController controllers.ItemController
}

func NewRouteItemController(itemController controllers.ItemController) ItemRouteController {
	return ItemRouteController{itemController}
}

func (ic *ItemRouteController) ItemRoute(rg *gin.RouterGroup) {
	router := rg.Group("items")
	router.GET("", ic.itemController.GetItems)
	router.GET("/:id", ic.itemController.GetItem)
	router.POST("", middleware.DeserializeUser(), middleware.CheckUserRole(), ic.itemController.CreateItem)
	router.PUT("/:id", ic.itemController.UpdateItem)
	router.DELETE("/:id", ic.itemController.DeleteItem)
	router.POST("/rating/:id", middleware.DeserializeUser(), ic.itemController.GiveRatingToItem)
	router.POST("/comment/:id", middleware.DeserializeUser(), ic.itemController.CommentItem)
}

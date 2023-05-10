package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gmkanat/Go-Shop/controllers"
	"github.com/gmkanat/Go-Shop/initializers"
	"github.com/gmkanat/Go-Shop/middleware"
	"github.com/gmkanat/Go-Shop/routes"
	"log"
	"net/http"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	ItemController      controllers.ItemController
	ItemRouteController routes.ItemRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	ItemController = controllers.NewItemController(initializers.DB)
	ItemRouteController = routes.NewRouteItemController(ItemController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	router := server.Group("/api")
	router.GET("/health-checker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})
	router.POST("order/:id/status", middleware.DeserializeUser(), middleware.CheckUserRole(ItemController.DB), ItemController.OrderStatus)
	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	ItemRouteController.ItemRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))
}

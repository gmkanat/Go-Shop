package main

import (
	"fmt"
	"github.com/gmkanat/Go-Shop/initializers"
	"github.com/gmkanat/Go-Shop/models"
	"log"
)

func init() {
	config, err := initializers.LoadConfig("/Users/User/go/src/github.com/gmkanat/Go-Shop")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	//if initializers.DB.Migrator().HasTable(&models.User{}) {
	//	initializers.DB.Migrator().DropTable(&models.User{})
	//}
	//if initializers.DB.Migrator().HasTable(&models.Item{}) {
	//	initializers.DB.Migrator().DropTable(&models.User{})
	//}
	//if initializers.DB.Migrator().HasTable(&models.ItemRating{}) {
	//	initializers.DB.Migrator().DropTable(&models.User{})
	//}
	initializers.DB.AutoMigrate(&models.User{}, models.UserRole{}, &models.Item{}, &models.ItemRating{}, &models.ItemComment{})
	fmt.Println("? Migration complete")
}

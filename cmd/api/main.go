package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nkyizbay/shop-project/internal/auth"
	"github.com/nkyizbay/shop-project/internal/item"
	"github.com/nkyizbay/shop-project/internal/user"
	"github.com/nkyizbay/shop-project/pkg/database"
	"github.com/spf13/viper"
)

func main() {
	router := gin.Default()

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	router.Use(auth.TokenMiddleware())

	fmt.Println(viper.Get("POSTGRES_DB")) // shop

	jwt := viper.GetString("SHOP_GO_JWTKEY")

	connectionPool, err := database.Setup()
	if err != nil {
		log.Fatal(err)
	}
	database.Migrate()

	// Item
	itemRepo := item.NewItemRepository(connectionPool)
	itemService := item.NewItemService(itemRepo)
	item.Handler(router, itemService)

	// USER
	userRepository := user.NewRepository(connectionPool)
	userService := user.NewUserService(userRepository)
	user.NewHandler(router, userService, jwt)

	router.Run(":8080")
}

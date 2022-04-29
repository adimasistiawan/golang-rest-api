package main

import (
	"bwastartup/handler"
	"bwastartup/user"
	"log"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// https://github.com/go-sql-driver/mysql
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/register", userHandler.Register)
	api.POST("/login", userHandler.Login)
	router.Run()
}

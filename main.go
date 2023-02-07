package main

import (
	"dairy_service/controller"
	"dairy_service/database"
	"dairy_service/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

// In this main function, the environment variables are loaded with the loadEnv function and the connection is established
// with the database using the loadDatabase function. If the connection was opened successfully, the AutoMigrate() function is
// called to create the relevant tables and columns for the User and Entry structs (if they don’t already exist).
// In addition to loading the environment variables and database connection, i am also creating a Gin router and declaring
// two routes for registration and login respectively using the serveApplication function.
func main() {
	loadEnv()
	loadDatabase()
	serveApplication()
}

func loadDatabase() {
	database.Connect()
	database.Database.AutoMigrate(&model.User{})
	database.Database.AutoMigrate(&model.Entry{})
}

func loadEnv() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func serveApplication() {
	router := gin.Default()

	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/register", controller.Register)
	publicRoutes.POST("/login", controller.Login)

	router.Run(":8000")
	fmt.Println("Server running on port 8000")
}

package main

import (
	"fmt"
	"log"

	"github.com/Iretoms/hng-task-two/config"
	"github.com/Iretoms/hng-task-two/model"
	"github.com/Iretoms/hng-task-two/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	loadDatabase()
	serveApp()
}

func loadDatabase() {
	config.Connect()
	config.Database.AutoMigrate(&model.User{})
	config.Database.AutoMigrate(&model.Organisation{})
}

func serveApp() {
	router := gin.Default()

	apiRoutes := router.Group("/api")
	authRoutes := router.Group("/auth")

	routes.RegisterRoute(authRoutes)
	routes.LoginRoute(authRoutes)
	routes.UserRoutes(apiRoutes)
	routes.OrganisationRoutes(apiRoutes)

	router.Run(":8080")
	fmt.Println("Server running on port 8080")
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/Iretoms/hng-task-two/config"
	"github.com/Iretoms/hng-task-two/middleware"
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
	err := config.Database.AutoMigrate(&model.User{}, &model.Organisation{})
	if err != nil {
		log.Fatalf("Could not migrate the database: %v", err)
	}
}

func serveApp() {
	router := gin.Default()

	protectedRoutes := router.Group("/api")
	publicRoutes := router.Group("/auth")

	routes.RegisterRoute(publicRoutes)
	routes.LoginRoute(publicRoutes)

	protectedRoutes.Use(middleware.JWTAuthMiddleware())
	routes.UserRoutes(protectedRoutes)
	routes.OrganisationRoutes(protectedRoutes)

	router.Run(":8080")
	fmt.Println("Server running on port 8080")
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

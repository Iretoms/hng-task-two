package routes

import "github.com/gin-gonic/gin"

func UserRoutes(r *gin.RouterGroup) {
	r.GET("/users/:id")
}

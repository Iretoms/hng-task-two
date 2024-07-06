package routes

import (
	"github.com/Iretoms/hng-task-two/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	r.GET("/users/:id", controllers.GetUser())
}

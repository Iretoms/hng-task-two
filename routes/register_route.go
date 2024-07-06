package routes

import (
	"github.com/Iretoms/hng-task-two/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(r *gin.RouterGroup) {
	r.POST("/register", controllers.Register())
}

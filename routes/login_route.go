package routes

import (
	"github.com/Iretoms/hng-task-two/controllers"
	"github.com/gin-gonic/gin"
)

func LoginRoute(r *gin.RouterGroup) {
	r.POST("/login", controllers.Login())
}

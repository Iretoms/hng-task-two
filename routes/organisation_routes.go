package routes

import (
	"github.com/Iretoms/hng-task-two/controllers"
	"github.com/gin-gonic/gin"
)

func OrganisationRoutes(r *gin.RouterGroup) {
	r.POST("/organisations", controllers.CreateOrganisation())
	r.POST("/organisations/:orgId/users", controllers.AddUserToOrganisation())
	r.GET("/organisations/:orgId", controllers.GetSingleOrganisation())
	r.GET("/organisations", controllers.GetUserOrganisations())
}

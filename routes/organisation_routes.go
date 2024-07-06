package routes

import "github.com/gin-gonic/gin"

func OrganisationRoutes(r *gin.RouterGroup) {
	r.POST("/organisations")
	r.POST("/organisations/:orgId/users")
	r.GET("/organisations/:orgId")
	r.GET("/organisations")
}
package controllers

import (
	"net/http"

	"github.com/Iretoms/hng-task-two/model"
	"github.com/Iretoms/hng-task-two/responses"
	"github.com/gin-gonic/gin"
)

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id")

		user, _ := model.FindUserById(userId)

		c.JSON(http.StatusOK, gin.H{
			"status":  "Success",
			"message": "Record gotten successfully",
			"data": responses.UserRes{
				UserID:    user.UserID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Phone:     user.Phone,
			},
		})
	}
}

package controllers

import (
	"net/http"
	"strings"

	"github.com/Iretoms/hng-task-two/config"
	"github.com/Iretoms/hng-task-two/helper"
	"github.com/Iretoms/hng-task-two/model"
	"github.com/Iretoms/hng-task-two/responses"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserOrganisations() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := helper.CurrentUser(c)

		var organisations []responses.OrganisationRes
		for _, org := range user.Organisations {
			orgRes := responses.OrganisationRes{
				OrgID:       org.OrgID,
				Name:        org.Name,
				Description: org.Description,
			}
			organisations = append(organisations, orgRes)
		}

		c.JSON(http.StatusOK, responses.SuccessResponse{
			Status:  "Success",
			Message: "User Organisations fetch successful",
			Data: responses.Data{
				Organisations: organisations,
			},
		})
	}
}

func GetSingleOrganisation() gin.HandlerFunc {
	return func(c *gin.Context) {
		orgId := c.Param("orgId")

		organisation, _ := model.FindOrganisationById(orgId)

		c.JSON(http.StatusOK, gin.H{
			"status":  "Success",
			"message": "Organisation gotten successfully",
			"data": responses.OrganisationRes{
				OrgID:       organisation.OrgID,
				Name:        organisation.Name,
				Description: organisation.Description,
			},
		})

	}
}

func CreateOrganisation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input model.OrgInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":     "Bad request",
				"message":    "Client error",
				"statusCode": "400",
			})
			return
		}

		user, err := helper.CurrentUser(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":     "Bad request",
				"message":    "Client error",
				"statusCode": "400",
			})
			return
		}

		uuid := uuid.New()
		uuidString := strings.Replace(uuid.String(), "-", "", -1)

		org := model.Organisation{
			OrgID:       uuidString,
			Name:        input.Name,
			Description: input.Description,
			Users:       []*model.User{&user},
		}

		newOrg, err := org.Save()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":     "Bad request",
				"message":    "Client error",
				"statusCode": "400",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  "Success",
			"message": "Organisation created successfully",
			"data": responses.OrganisationRes{
				OrgID:       newOrg.OrgID,
				Name:        newOrg.Name,
				Description: newOrg.Description,
			}})
	}
}

func AddUserToOrganisation() gin.HandlerFunc {
	return func(c *gin.Context) {
		orgId := c.Param("orgId")

		var userInput model.AddUserInput
		if err := c.ShouldBindJSON(&userInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}

		organisation, _ := model.FindOrganisationById(orgId)
		user, _ := model.FindUserById(userInput.UserID)

		organisation.Users = append(organisation.Users, &user)

		if err := config.Database.Save(&organisation).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to organisation: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "Success",
			"message": "User added to organisation successfully",
		})
	}
}

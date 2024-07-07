package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Iretoms/hng-task-two/helper"
	"github.com/Iretoms/hng-task-two/model"
	"github.com/Iretoms/hng-task-two/responses"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input model.RegisterInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(
				http.StatusUnprocessableEntity,
				responses.ErrorResponse{
					Errors: helper.FormatValidationError(err),
				})
			return
		}

		err := validate.Struct(input)
		if err != nil {
			c.JSON(
				http.StatusUnprocessableEntity,
				responses.ErrorResponse{
					Errors: helper.FormatValidationError(err),
				})
			return
		}

		uuid1 := uuid.New()
		uuidString1 := strings.Replace(uuid1.String(), "-", "", -1)

		uuid2 := uuid.New()
		uuidString2 := strings.Replace(uuid2.String(), "-", "", -1)

		user := model.User{
			UserID:    uuidString1,
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Password:  input.Password,
			Phone:     input.Phone,
			Organisations: []*model.Organisation{{
				OrgID:       uuidString2,
				Name:        fmt.Sprintf("%v's Organisation", input.FirstName),
				Description: "",
			}},
		}

		savedUser, err := user.Save()

		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				responses.ErrorResponse{
					Status:     "Bad request",
					Message:    "Registration Unsuccessful",
					StatusCode: http.StatusBadRequest,
				})
			return
		}

		jwt, err := helper.GenerateJWT(user)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				responses.ErrorResponse{
					Status:     "Bad request",
					Message:    "Registration Unsuccessful",
					StatusCode: http.StatusBadRequest,
				})
			return
		}

		c.JSON(
			http.StatusCreated,
			responses.SuccessResponse{
				Status:  "success",
				Message: "Registration successful",
				Data: responses.Data{
					AccessToken: jwt,
					User: &responses.UserRes{
						UserID:    savedUser.UserID,
						FirstName: savedUser.FirstName,
						LastName:  savedUser.LastName,
						Email:     savedUser.Email,
						Phone:     savedUser.Phone},
				},
			})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input model.LoginInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(
				http.StatusUnprocessableEntity,
				responses.ErrorResponse{
					Errors: helper.FormatValidationError(err),
				})
			return
		}

		err := validate.Struct(input)
		if err != nil {
			c.JSON(
				http.StatusUnprocessableEntity,
				responses.ErrorResponse{
					Errors: helper.FormatValidationError(err),
				})
			return
		}

		user, err := model.FindUserByEmail(input.Email)

		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				responses.ErrorResponse{
					Status:     "Bad request",
					Message:    "Authentication failed",
					StatusCode: http.StatusUnauthorized,
				})
			return
		}

		err = user.ValidatePassword(input.Password)

		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				responses.ErrorResponse{
					Status:     "Bad request",
					Message:    "Authentication failed",
					StatusCode: http.StatusUnauthorized,
				})
			return
		}

		jwt, err := helper.GenerateJWT(user)
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				responses.ErrorResponse{
					Status:     "Bad request",
					Message:    "Authentication failed",
					StatusCode: http.StatusUnauthorized,
				})
			return
		}

		c.JSON(
			http.StatusOK,
			responses.SuccessResponse{
				Status:  "success",
				Message: "Login successful",
				Data: responses.Data{
					AccessToken: jwt,
					User: &responses.UserRes{
						UserID:    user.UserID,
						FirstName: user.FirstName,
						LastName:  user.LastName,
						Email:     user.Email,
						Phone:     user.Phone,
					}}})

	}
}

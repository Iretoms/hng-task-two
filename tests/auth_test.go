package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Iretoms/hng-task-two/config"
	"github.com/Iretoms/hng-task-two/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	config.Connect()

	publicRoutes := router.Group("/auth")
	routes.RegisterRoute(publicRoutes)
	routes.LoginRoute(publicRoutes)

	return router
}
func TestUserAuth(t *testing.T) {
	router := setupRouter()

	t.Run("It Should Register User Successfully with Default Organisation", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"firstName": "John",
			"lastName":  "Doe",
			"email":     "john.doe@example.com",
			"password":  "password123",
			"phone":     "1234567890",
		})

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		user := data["user"].(map[string]interface{})
		organisation := user["organisation"].([]interface{})[0].(map[string]interface{})

		assert.Equal(t, "John", user["firstName"])
		assert.Equal(t, "Doe", user["lastName"])
		assert.Equal(t, "john.doe@example.com", user["email"])
		assert.Equal(t, "1234567890", user["phone"])
		assert.Equal(t, "John's Organisation", organisation["name"])
		assert.NotEmpty(t, data["accessToken"])
	})

	t.Run("It Should Log the user in successfully", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"email":    "john.doe@example.com",
			"password": "password123",
		})

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		user := data["user"].(map[string]interface{})

		assert.Equal(t, "john.doe@example.com", user["email"])
		assert.NotEmpty(t, data["accessToken"])
	})

	t.Run("It Should Fail If Required Fields Are Missing", func(t *testing.T) {
		requiredFields := []string{"firstName", "lastName", "email", "password"}

		for _, field := range requiredFields {
			requestBody := map[string]string{
				"firstName": "John",
				"lastName":  "Doe",
				"email":     "john.doe@example.com",
				"password":  "password123",
				"phone":     "08028487596",
			}
			delete(requestBody, field)

			body, _ := json.Marshal(requestBody)
			req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			errors := response["errors"].([]interface{})
			assert.NotEmpty(t, errors)

			var errorFieldFound bool
			for _, err := range errors {
				errorMap := err.(map[string]interface{})
				if errorMap["field"] == field {
					errorFieldFound = true
					assert.Equal(t, "Invalid input", errorMap["message"])
					break
				}
			}

			assert.True(t, errorFieldFound, fmt.Sprintf("Expected error for field %s", field))
		}
	})

	t.Run("It Should Fail if there’s Duplicate Email or UserID", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"firstName": "Jane",
			"lastName":  "Doe",
			"email":     "john.doe@example.com",
			"password":  "password123",
			"phone":     "09034087736",
		})

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response["message"], "Email already exists")
	})

}

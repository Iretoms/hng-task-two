package tests

import (
	"bytes"
	"encoding/json"
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

	t.Run("It_Should_Register_User_Successfully_with_Default_Organisation", func(t *testing.T) {
		requestBody := map[string]string{
			"firstName": "John",
			"lastName":  "Doe",
			"email":     "john.do123@example.com",
			"password":  "password123",
			"phone":     "1234567890",
		}
		body, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "success", response["status"])
		assert.Equal(t, "Registration successful", response["message"])

		data, dataExists := response["data"].(map[string]interface{})
		assert.True(t, dataExists, "Response data should exist")

		user, userExists := data["user"].(map[string]interface{})
		assert.True(t, userExists, "User data should exist")

		assert.Equal(t, "John", user["firstName"])
		assert.Equal(t, "Doe", user["lastName"])
		assert.Equal(t, "john.do123@example.com", user["email"])

		organisations, orgExists := user["organisations"].([]interface{})
		if !orgExists || organisations == nil || len(organisations) == 0 {
			t.Fatalf("Organisation field is missing or empty. Full response: %v", response)
		}

		organisation := organisations[0].(map[string]interface{})
		assert.Equal(t, "John's Organisation", organisation["name"])
	})

	t.Run("It_Should_Log_the_user_in_successfully", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"email":    "john.do123@example.com",
			"password": "password123",
		})

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Error unmarshalling response")

		data, dataExists := response["data"].(map[string]interface{})
		assert.True(t, dataExists, "Response data should exist")

		user, userExists := data["user"].(map[string]interface{})
		assert.True(t, userExists, "User data should exist")

		assert.Equal(t, "john.do123@example.com", user["email"])
		assert.NotEmpty(t, data["accessToken"])
	})


	t.Run("It_Should_Fail_If_Required_Fields_Are_Missing", func(t *testing.T) {
		requiredFields := []string{"FirstName", "LastName", "Email", "Password"}

		for _, field := range requiredFields {
			requestBody := map[string]string{
				"FirstName": "John",
				"LastName":  "Doe",
				"Email":     "john.do@example.com",
				"Password":  "password123",
				"Phone":     "1234567890",
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

			// Log the entire response for debugging purposes
			t.Logf("Response body: %v", response)

			errors, ok := response["errors"].([]interface{})
			if !ok {
				t.Fatalf("Expected errors in response. Full response: %v", response)
			}

			// Check if the expected error message for the missing field is present
			found := false
			for _, err := range errors {
				errMap := err.(map[string]interface{})
				if errMap["field"] == field {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected error for field %s", field)
		}
	})

	t.Run("It Should Fail if there’s Duplicate Email or UserID", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]string{
			"firstName": "Jane",
			"lastName":  "Doe",
			"email":     "john.do123@example.com",
			"password":  "password123",
			"phone":     "09034087736",
		})

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response["message"], "Registration unsuccessful")
	})

}

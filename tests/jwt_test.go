package tests

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/Iretoms/hng-task-two/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var privateKey []byte

func GenerateJWT(user model.User) (string, error) {
	tokenTTLMinutes, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))

	now := time.Now()
	iat := now.Unix()
	exp := now.Add(time.Duration(tokenTTLMinutes) * time.Minute).Unix()

	fmt.Printf("iat: %v, exp: %v, difference: %v minutes\n", iat, exp, tokenTTLMinutes)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.UserID,
		"iat": iat,
		"exp": exp,
	})
	return token.SignedString(privateKey)
}

func TestGenerateJWT(t *testing.T) {
	loadEnv()

	if len(privateKey) == 0 {
		t.Fatal("JWT_SECRET_KEY environment variable is not set")
	}

	user := model.User{UserID: "85d7b1717583401083e5d5c1c85edc8c"}

	tokenString, err := GenerateJWT(user)
	assert.NoError(t, err, "Token generation should not produce an error")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	assert.NoError(t, err, "Token parsing should not produce an error")

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		assert.Equal(t, user.UserID, claims["id"], "UserID should match")

		expirationAt := int64(claims["exp"].(float64))

		expectedExpiration := time.Now().Add(15 * time.Minute)

		assert.WithinDuration(t, time.Unix(expirationAt, 0), expectedExpiration, 5*time.Second, "Expiration time should be close to 15 minutes from now")

	} else {
		t.Errorf("Token claims are invalid or token is not valid")
	}
}

func loadEnv() {
	projectRoot, err := filepath.Abs("..")
	if err != nil {
		log.Fatalf("Error determining project root: %v", err)
	}

	envPath := filepath.Join(projectRoot, ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	privateKey = []byte(os.Getenv("JWT_SECRET_KEY"))
}

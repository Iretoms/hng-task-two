package helper

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Iretoms/hng-task-two/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var privateKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateJWT(user model.User) (string, error) {
	tokenTTLMinutes, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	
	now := time.Now()
	iat := now.Unix()
	exp := now.Add(time.Duration(tokenTTLMinutes) * time.Minute).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.UserID,
		"iat": iat,
		"exp": exp,
	})
	return token.SignedString(privateKey)
}

func ValidateJWT(c *gin.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

func CurrentUser(c *gin.Context) (model.User, error) {
	err := ValidateJWT(c)
	if err != nil {
		return model.User{}, err
	}

	token, err := getToken(c)
	if err != nil {
		return model.User{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.User{}, errors.New("invalid token claims")
	}

	userId, ok := claims["id"].(string)
	if !ok {
		return model.User{}, errors.New("user ID not found in token claims")
	}

	user, err := model.FindUserById(userId)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func getToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

func getTokenFromRequest(context *gin.Context) string {
	bearerToken := context.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

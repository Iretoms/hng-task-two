package model

import (
	"html"
	"strings"

	"github.com/Iretoms/hng-task-two/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	UserID        string          `gorm:"primaryKey" json:"userId" validate:"required,uuid"`
	FirstName     string          `gorm:"size:255;not null" json:"firstName" validate:"required"`
	LastName      string          `gorm:"size:255;not null" json:"lastName" validate:"required"`
	Email         string          `gorm:"size:255;not null;unique" json:"email" validate:"required,email"`
	Password      string          `gorm:"size:255;not null" json:"-" validate:"required"`
	Phone         string          `json:"phone"`
	Organisations []*Organisation `gorm:"many2many:user_organisations;"`
}

type AddUserInput struct {
	UserID string `json:"userId" binding:"required"`
}

func (user *User) Save() (*User, error) {
	err := config.Database.Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) BeforeSave(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(passwordHash)
	user.FirstName = html.EscapeString(strings.TrimSpace(user.FirstName))
	user.LastName = html.EscapeString(strings.TrimSpace(user.LastName))
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	return nil
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func FindUserByEmail(email string) (User, error) {
	var user User
	err := config.Database.Where("email=?", email).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func FindUserById(id string) (User, error) {
	var user User

	err := config.Database.Preload("Organisations").Where("user_id=?", id).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

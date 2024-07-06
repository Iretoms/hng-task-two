package model

type Organisation struct {
	OrgID       string  `gorm:"primaryKey" json:"orgId"`
	Name        string  `gorm:"size:255;not null" json:"name" validate:"required"`
	Description string  `json:"description"`
	Users       []*User `gorm:"many2many:user_organisations;"`
}

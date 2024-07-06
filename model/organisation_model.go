package model

import "github.com/Iretoms/hng-task-two/config"

type Organisation struct {
	OrgID       string  `gorm:"primaryKey" json:"orgId"`
	Name        string  `gorm:"size:255;not null" json:"name" validate:"required"`
	Description string  `json:"description"`
	Users       []*User `gorm:"many2many:user_organisations;"`
}

type OrgInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (org *Organisation) Save() (*Organisation, error) {
	err := config.Database.Create(&org).Error
	if err != nil {
		return &Organisation{}, err
	}
	return org, nil
}

func FindOrganisationById(id string) (Organisation, error) {
	var organisation Organisation

	err := config.Database.Preload("Users").Where("org_id=?", id).Find(&organisation).Error
	if err != nil {
		return Organisation{}, err
	}
	return organisation, nil
}

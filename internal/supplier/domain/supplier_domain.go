package domain

import (
	"inv/pkg/domain"
)

type Supplier struct {
	domain.BaseModel
	Name     string `gorm:"type:varchar(100);not null"`
	Code     string `gorm:"type:varchar(50);unique;not null"`
	Contact  string `gorm:"type:varchar(100)"`
	Email    string `gorm:"type:varchar(255)"`
	Phone    string `gorm:"type:varchar(50)"`
	IsActive bool   `gorm:"default:true"`
}

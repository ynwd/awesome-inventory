package domain

import (
	"inv/pkg/domain"
)

type Category struct {
	domain.BaseModel
	Name        string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:text"`
}

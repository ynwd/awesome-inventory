package domain

import (
	"inv/pkg/domain"
)

type Inventory struct {
	domain.BaseModel
	ProductID string `gorm:"type:varchar(26);index;not null"`
	Quantity  int    `gorm:"not null;default:0"`
	Location  string `gorm:"type:varchar(100)"`
}

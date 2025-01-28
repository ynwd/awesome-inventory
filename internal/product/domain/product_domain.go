package domain

import (
	"inv/pkg/domain"
)

type Product struct {
	domain.BaseModel
	Name        string  `gorm:"type:varchar(100);not null"`
	SKU         string  `gorm:"type:varchar(50);unique;not null"`
	SupplierID  string  `gorm:"type:varchar(26);index;not null"`
	Description string  `gorm:"type:text"`
	UnitPrice   float64 `gorm:"type:decimal(10,2);not null"`
}

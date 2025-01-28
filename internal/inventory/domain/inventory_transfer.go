package domain

import (
	"inv/pkg/domain"
	"time"
)

type MovementType string

const (
	TransferIn  MovementType = "TRANSFER_IN"
	TransferOut MovementType = "TRANSFER_OUT"
)

type InventoryTransfer struct {
	domain.BaseModel
	InventoryID    string       `gorm:"type:varchar(26);index;not null"`
	Type           MovementType `gorm:"type:varchar(20);not null"`
	Quantity       int          `gorm:"not null"`
	SourceLocation string       `gorm:"type:varchar(100);not null"`
	DestLocation   string       `gorm:"type:varchar(100);not null"`
	TransferDate   time.Time    `gorm:"not null"`
	Notes          string       `gorm:"type:text"`
}

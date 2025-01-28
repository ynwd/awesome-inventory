package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string         `gorm:"type:varchar(26);primary_key"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ulid.Make().String()
	}
	return nil
}

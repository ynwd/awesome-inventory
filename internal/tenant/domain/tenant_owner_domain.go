package domain

import (
	"errors"
	"inv/pkg/domain"

	"golang.org/x/crypto/bcrypt"
)

type TenantOwner struct {
	domain.BaseModel
	Email    string   `gorm:"type:varchar(255);unique;not null"`
	Name     string   `gorm:"type:varchar(100);not null"`
	Password string   `gorm:"type:varchar(255);not null"`
	Tenants  []Tenant `gorm:"foreignKey:TenantOwnerID"`
}

func (o *TenantOwner) Validate() error {
	if o.Email == "" {
		return errors.New("email is required")
	}
	if o.Name == "" {
		return errors.New("name is required")
	}
	if o.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (o *TenantOwner) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(o.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	o.Password = string(hashedPassword)
	return nil
}

func (o *TenantOwner) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(o.Password), []byte(password))
	return err == nil
}

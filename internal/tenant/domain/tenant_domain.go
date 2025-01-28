package domain

import (
	"errors"
	"fmt"
	"inv/pkg/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TenantStatus string

const (
	TenantActive   TenantStatus = "ACTIVE"
	TenantInactive TenantStatus = "INACTIVE"
)

type Tenant struct {
	domain.BaseModel
	TenantOwnerID string       `json:"tenant_owner_id" gorm:"type:uuid;not null"`
	DatabaseName  string       `json:"database_name" gorm:"type:varchar(100);not null"`
	DatabaseHost  string       `json:"database_host" gorm:"type:varchar(255);not null"`
	DatabasePort  string       `json:"database_port" gorm:"type:varchar(10);not null"`
	DatabaseUser  string       `json:"database_user" gorm:"type:varchar(100);not null"`
	DatabasePass  string       `json:"database_pass" gorm:"type:varchar(255);not null"`
	Status        TenantStatus `json:"status" gorm:"type:varchar(20);not null"`
	Owner         *TenantOwner `gorm:"foreignKey:TenantOwnerID;references:ID"`
}

func (t *Tenant) Validate() error {
	if t.DatabaseName == "" {
		return errors.New("database name is required")
	}
	if t.DatabaseHost == "" {
		return errors.New("database host is required")
	}
	if t.DatabasePort == "" {
		return errors.New("database port is required")
	}
	if t.DatabaseUser == "" {
		return errors.New("database user is required")
	}
	if t.DatabasePass == "" {
		return errors.New("database password is required")
	}
	return nil
}

func (t *Tenant) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		t.DatabaseHost,
		t.DatabasePort,
		t.DatabaseUser,
		t.DatabasePass,
		t.DatabaseName,
	)
}

func (t *Tenant) IsActivated() bool {
	return t.Status == TenantActive
}

func (t *Tenant) CreateDatabase(db *gorm.DB) error {
	// Create connection string to connect to postgres database
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		t.DatabaseHost,
		t.DatabasePort,
		t.DatabaseUser,
		t.DatabasePass,
	)

	// Open new connection
	postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Check if database already exists
	var exists bool
	checkSQL := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = '%s')", t.DatabaseName)
	if err := postgresDB.Raw(checkSQL).Scan(&exists).Error; err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if exists {
		return fmt.Errorf("database %s already exists", t.DatabaseName)
	}

	// Create the new database
	createDBSQL := fmt.Sprintf("CREATE DATABASE %s", t.DatabaseName)
	if err := postgresDB.Exec(createDBSQL).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}

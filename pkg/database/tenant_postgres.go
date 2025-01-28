package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TenantConfig extends Config with tenant-specific database name
type TenantConfig struct {
	*Config
	TenantDBName string
}

// NewTenantConnection creates a connection to a tenant-specific database
func NewTenantConnection(tenantConfig *TenantConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		tenantConfig.Host,
		tenantConfig.Port,
		tenantConfig.User,
		tenantConfig.Password,
		tenantConfig.TenantDBName,
		tenantConfig.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tenant database: %w", err)
	}

	return db, nil
}

// GetTenantConnection retrieves a tenant connection using the tenant identifier
func GetTenantConnection(mainDB *gorm.DB, tenantID string, baseConfig *Config) (*gorm.DB, error) {
	// You might want to adjust this query based on your tenant table structure
	var tenantDBName string
	err := mainDB.Table("tenants").
		Select("database_name").
		Where("id = ?", tenantID).
		First(&tenantDBName).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant database name: %w", err)
	}

	tenantConfig := &TenantConfig{
		Config:       baseConfig,
		TenantDBName: tenantDBName,
	}

	return NewTenantConnection(tenantConfig)
}

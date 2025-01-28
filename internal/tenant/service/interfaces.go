package service

import (
	"context"
	"inv/internal/tenant/domain"

	"gorm.io/gorm"
)

type TenantService interface {
	CreateTenant(ctx context.Context, tenant *domain.Tenant) error
	UpdateTenant(ctx context.Context, tenant *domain.Tenant) error
	DeleteTenant(ctx context.Context, id string) error
	GetTenant(ctx context.Context, id string) (*domain.Tenant, error)
	GetTenantByOwner(ctx context.Context, ownerID string) (*domain.Tenant, error)
	ListTenants(ctx context.Context) ([]*domain.Tenant, error)
	UpdateTenantStatus(ctx context.Context, id string, status domain.TenantStatus) error
}

type TenantOwnerService interface {
	CreateOwner(ctx context.Context, owner *domain.TenantOwner) error
	UpdateOwner(ctx context.Context, owner *domain.TenantOwner) error
	DeleteOwner(ctx context.Context, id string) error
	GetOwner(ctx context.Context, id string) (*domain.TenantOwner, error)
	GetOwnerByEmail(ctx context.Context, email string) (*domain.TenantOwner, error)
	ListOwners(ctx context.Context, page, pageSize int) ([]*domain.TenantOwner, error)
}

type TenantDatabaseService interface {
	CreateDatabase(ctx context.Context, tenant *domain.Tenant) error
	MigrateSchema(ctx context.Context, tenant *domain.Tenant) error
	DropDatabase(ctx context.Context, tenant *domain.Tenant) error
	GetConnection(ctx context.Context, tenant *domain.Tenant) (*gorm.DB, error)
	CloseConnection(ctx context.Context, tenant *domain.Tenant) error
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

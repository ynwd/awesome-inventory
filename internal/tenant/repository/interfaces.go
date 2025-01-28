package repository

import (
	"context"
	"errors"
	"inv/internal/tenant/domain"
)

var (
	ErrTenantNotFound = errors.New("tenant not found")
	ErrTenantExists   = errors.New("tenant already exists")
)

type TenantOwnerRepository interface {
	Create(owner *domain.TenantOwner) error
	Update(owner *domain.TenantOwner) error
	Delete(id string) error
	FindByID(id string) (*domain.TenantOwner, error)
	FindByEmail(email string) (*domain.TenantOwner, error)
	List(page, pageSize int) ([]*domain.TenantOwner, error)
}

type TenantRepository interface {
	Create(ctx context.Context, tenant *domain.Tenant) error
	Update(ctx context.Context, tenant *domain.Tenant) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*domain.Tenant, error)
	FindByOwnerID(ctx context.Context, ownerID string) (*domain.Tenant, error)
	FindAll(ctx context.Context) ([]*domain.Tenant, error)
	FindByDatabaseName(ctx context.Context, dbName string) (*domain.Tenant, error)
	UpdateStatus(ctx context.Context, id string, status domain.TenantStatus) error
}

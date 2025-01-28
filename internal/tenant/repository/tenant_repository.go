package repository

import (
	"context"
	"inv/internal/tenant/domain"

	"gorm.io/gorm"
)

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	result := r.db.WithContext(ctx).Create(tenant)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *tenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	result := r.db.WithContext(ctx).Save(tenant)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTenantNotFound
	}
	return nil
}

func (r *tenantRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.Tenant{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTenantNotFound
	}
	return nil
}

func (r *tenantRepository) FindByID(ctx context.Context, id string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	result := r.db.WithContext(ctx).Preload("Owner").First(&tenant, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrTenantNotFound
		}
		return nil, result.Error
	}
	return &tenant, nil
}

func (r *tenantRepository) FindByOwnerID(ctx context.Context, ownerID string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	result := r.db.WithContext(ctx).Preload("Owner").First(&tenant, "tenant_owner_id = ?", ownerID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrTenantNotFound
		}
		return nil, result.Error
	}
	return &tenant, nil
}

func (r *tenantRepository) FindAll(ctx context.Context) ([]*domain.Tenant, error) {
	var tenants []*domain.Tenant
	result := r.db.WithContext(ctx).Preload("Owner").Find(&tenants)
	if result.Error != nil {
		return nil, result.Error
	}
	return tenants, nil
}

func (r *tenantRepository) FindByDatabaseName(ctx context.Context, dbName string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	result := r.db.WithContext(ctx).Preload("Owner").First(&tenant, "database_name = ?", dbName)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrTenantNotFound
		}
		return nil, result.Error
	}
	return &tenant, nil
}

func (r *tenantRepository) UpdateStatus(ctx context.Context, id string, status domain.TenantStatus) error {
	result := r.db.WithContext(ctx).Model(&domain.Tenant{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTenantNotFound
	}
	return nil
}

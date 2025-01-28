package repository

import (
	"inv/internal/tenant/domain"

	"gorm.io/gorm"
)

type tenantOwnerRepository struct {
	db *gorm.DB
}

func NewTenantOwnerRepository(db *gorm.DB) TenantOwnerRepository {
	return &tenantOwnerRepository{db: db}
}

func (r *tenantOwnerRepository) Create(owner *domain.TenantOwner) error {
	// Check if email already exists
	exists := &domain.TenantOwner{}
	err := r.db.Where("email = ?", owner.Email).First(exists).Error
	if err == nil {
		return ErrTenantExists
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	return r.db.Create(owner).Error
}

func (r *tenantOwnerRepository) Update(owner *domain.TenantOwner) error {
	return r.db.Save(owner).Error
}

func (r *tenantOwnerRepository) Delete(id string) error {
	return r.db.Delete(&domain.TenantOwner{}, "id = ?", id).Error
}

func (r *tenantOwnerRepository) FindByID(id string) (*domain.TenantOwner, error) {
	var owner domain.TenantOwner
	err := r.db.Preload("Tenants").First(&owner, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return &owner, nil
}

func (r *tenantOwnerRepository) FindByEmail(email string) (*domain.TenantOwner, error) {
	var owner domain.TenantOwner
	err := r.db.Preload("Tenants").First(&owner, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return &owner, nil
}

func (r *tenantOwnerRepository) List(page, pageSize int) ([]*domain.TenantOwner, error) {
	var owners []*domain.TenantOwner
	offset := (page - 1) * pageSize
	err := r.db.Preload("Tenants").Offset(offset).Limit(pageSize).Find(&owners).Error
	if err != nil {
		return nil, err
	}
	return owners, nil
}

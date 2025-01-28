package service

import (
	"context"
	"inv/internal/tenant/domain"
	"inv/internal/tenant/repository"
)

type tenantOwnerService struct {
	ownerRepo repository.TenantOwnerRepository
}

func NewTenantOwnerService(repo repository.TenantOwnerRepository) TenantOwnerService {
	return &tenantOwnerService{ownerRepo: repo}
}

func (s *tenantOwnerService) CreateOwner(ctx context.Context, owner *domain.TenantOwner) error {
	if owner == nil {
		return ErrInvalidInput
	}

	if err := owner.Validate(); err != nil {
		return err
	}

	// Check if email is already in use
	existing, err := s.ownerRepo.FindByEmail(owner.Email)
	if err == nil && existing != nil {
		return ErrEmailExists
	}

	// Hash the password before saving
	if err := owner.HashPassword(); err != nil {
		return err
	}

	return s.ownerRepo.Create(owner)
}

func (s *tenantOwnerService) UpdateOwner(ctx context.Context, owner *domain.TenantOwner) error {
	if owner == nil {
		return ErrInvalidInput
	}

	if err := owner.Validate(); err != nil {
		return err
	}

	existing, err := s.ownerRepo.FindByID(owner.ID)
	if err != nil {
		return err
	}

	// Only hash password if it has changed
	if owner.Password != existing.Password {
		if err := owner.HashPassword(); err != nil {
			return err
		}
	}

	return s.ownerRepo.Update(owner)
}

func (s *tenantOwnerService) DeleteOwner(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidInput
	}
	return s.ownerRepo.Delete(id)
}

func (s *tenantOwnerService) GetOwner(ctx context.Context, id string) (*domain.TenantOwner, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}
	return s.ownerRepo.FindByID(id)
}

func (s *tenantOwnerService) GetOwnerByEmail(ctx context.Context, email string) (*domain.TenantOwner, error) {
	if email == "" {
		return nil, ErrInvalidInput
	}
	return s.ownerRepo.FindByEmail(email)
}

func (s *tenantOwnerService) ListOwners(ctx context.Context, page, pageSize int) ([]*domain.TenantOwner, error) {
	if page < 1 || pageSize < 1 {
		return nil, ErrInvalidInput
	}
	return s.ownerRepo.List(page, pageSize)
}

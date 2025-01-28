package service

import (
	"context"
	"fmt"
	"inv/internal/tenant/domain"
	"inv/internal/tenant/repository"
	"inv/pkg/utils"
)

type tenantService struct {
	tenantRepo repository.TenantRepository
	dbService  TenantDatabaseService
}

func NewTenantService(tenantRepo repository.TenantRepository, dbService TenantDatabaseService) TenantService {
	return &tenantService{
		tenantRepo: tenantRepo,
		dbService:  dbService,
	}
}

func (s *tenantService) CreateTenant(ctx context.Context, tenant *domain.Tenant) error {
	tenant.DatabaseName = fmt.Sprintf("tenant_%s", utils.GenerateULID())
	tenant.DatabaseUser = fmt.Sprintf("user_%s", utils.GenerateULID())
	tenant.DatabasePass = utils.GenerateSecurePassword()

	if err := s.dbService.CreateDatabase(ctx, tenant); err != nil {
		return err
	}

	if err := s.dbService.MigrateSchema(ctx, tenant); err != nil {
		_ = s.dbService.DropDatabase(ctx, tenant)
		return err
	}

	return s.tenantRepo.Create(ctx, tenant)
}

func (s *tenantService) UpdateTenant(ctx context.Context, tenant *domain.Tenant) error {
	return s.tenantRepo.Update(ctx, tenant)
}

func (s *tenantService) DeleteTenant(ctx context.Context, id string) error {
	return s.tenantRepo.Delete(ctx, id)
}

func (s *tenantService) GetTenant(ctx context.Context, id string) (*domain.Tenant, error) {
	return s.tenantRepo.FindByID(ctx, id)
}

func (s *tenantService) GetTenantByOwner(ctx context.Context, ownerID string) (*domain.Tenant, error) {
	return s.tenantRepo.FindByOwnerID(ctx, ownerID)
}

func (s *tenantService) ListTenants(ctx context.Context) ([]*domain.Tenant, error) {
	return s.tenantRepo.FindAll(ctx)
}

func (s *tenantService) UpdateTenantStatus(ctx context.Context, id string, status domain.TenantStatus) error {
	return s.tenantRepo.UpdateStatus(ctx, id, status)
}

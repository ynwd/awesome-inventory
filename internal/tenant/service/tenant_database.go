package service

import (
	"context"
	"fmt"
	ca "inv/internal/category/domain"
	it "inv/internal/inventory/domain"
	pr "inv/internal/product/domain"
	su "inv/internal/supplier/domain"
	"inv/internal/tenant/domain"
	us "inv/internal/user/domain"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type tenantDatabaseService struct {
	adminDB     *gorm.DB
	connections map[string]*gorm.DB
}

func NewTenantDatabaseService(adminDB *gorm.DB) TenantDatabaseService {
	return &tenantDatabaseService{
		adminDB:     adminDB,
		connections: make(map[string]*gorm.DB),
	}
}

// func (s *tenantDatabaseService) CreateDatabase(ctx context.Context, tenant *domain.Tenant) error {
// 	sql := fmt.Sprintf("CREATE DATABASE %s", tenant.DatabaseName)
// 	return s.adminDB.Exec(sql).Error
// }

func (s *tenantDatabaseService) CreateDatabase(ctx context.Context, tenant *domain.Tenant) error {
	createUserSQL := fmt.Sprintf(`CREATE USER "%s" WITH PASSWORD '%s'`,
		tenant.DatabaseUser,
		tenant.DatabasePass,
	)
	if err := s.adminDB.Exec(createUserSQL).Error; err != nil {
		s.adminDB.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, tenant.DatabaseName))
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Create database
	createDBSQL := fmt.Sprintf(`CREATE DATABASE "%s" WITH OWNER = '%s'`,
		tenant.DatabaseName,
		tenant.DatabaseUser, // Set the tenant user as the owner
	)
	if err := s.adminDB.Exec(createDBSQL).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	// Connect to the new database to set up schema
	tenantDSN := fmt.Sprintf(`host='%s' port=%s user='%s' password='%s' dbname=%s sslmode=disable`,
		tenant.DatabaseHost,
		tenant.DatabasePort,
		tenant.DatabaseUser,
		tenant.DatabasePass,
		tenant.DatabaseName,
	)

	tenantDB, err := gorm.Open(postgres.Open(tenantDSN), &gorm.Config{})
	if err != nil {
		s.adminDB.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, tenant.DatabaseName))
		s.adminDB.Exec(fmt.Sprintf(`DROP USER IF EXISTS "%s"`, tenant.DatabaseUser))
		return fmt.Errorf("failed to connect to tenant database: %w", err)
	}

	// Create schema and grant permissions
	if err := s.setupTenantSchema(tenantDB, tenant.DatabaseUser); err != nil {
		s.adminDB.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, tenant.DatabaseName))
		s.adminDB.Exec(fmt.Sprintf(`DROP USER IF EXISTS "%s"`, tenant.DatabaseUser))
		return fmt.Errorf("failed to set up tenant schema: %v", err)
	}

	return nil
}

func (s *tenantDatabaseService) setupTenantSchema(tenantDB *gorm.DB, dbUser string) any {
	statements := []string{
		// Ensure public schema exists and set permissions
		`CREATE SCHEMA IF NOT EXISTS public`,
		fmt.Sprintf(`GRANT ALL ON SCHEMA public TO "%s"`, dbUser),
		fmt.Sprintf(`ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO "%s"`, dbUser),
		fmt.Sprintf(`ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO "%s"`, dbUser),
		fmt.Sprintf(`ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO "%s"`, dbUser),
		fmt.Sprintf(`ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TYPES TO "%s"`, dbUser),
	}

	// Execute all schema setup statements
	for _, stmt := range statements {
		if err := tenantDB.Exec(stmt).Error; err != nil {
			return fmt.Errorf("failed to execute schema setup statement: %w", err)
		}
	}
	return nil
}

func (s *tenantDatabaseService) MigrateSchema(ctx context.Context, tenant *domain.Tenant) error {
	db, err := s.GetConnection(ctx, tenant)
	if err != nil {
		log.Default().Printf("MigrateSchema: Failed to get connection for tenant %s", tenant.ID)
		return err
	}

	// Define tenant-specific models
	return db.AutoMigrate(
		&us.User{},
		&ca.Category{},
		&pr.Product{},
		&su.Supplier{},
		&it.Inventory{},
		&it.InventoryTransfer{},
	)
}

func (s *tenantDatabaseService) GetConnection(ctx context.Context, tenant *domain.Tenant) (*gorm.DB, error) {
	if db, exists := s.connections[tenant.ID]; exists {
		return db, nil
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		tenant.DatabaseHost,
		tenant.DatabaseUser,
		tenant.DatabasePass,
		tenant.DatabaseName,
		tenant.DatabasePort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	s.connections[tenant.ID] = db
	return db, nil
}

func (s *tenantDatabaseService) DropDatabase(ctx context.Context, tenant *domain.Tenant) error {
	s.CloseConnection(ctx, tenant)
	sql := fmt.Sprintf("DROP DATABASE IF EXISTS %s", tenant.DatabaseName)
	return s.adminDB.Exec(sql).Error
}

func (s *tenantDatabaseService) CloseConnection(ctx context.Context, tenant *domain.Tenant) error {
	if db, exists := s.connections[tenant.ID]; exists {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		delete(s.connections, tenant.ID)
		return sqlDB.Close()
	}
	return nil
}

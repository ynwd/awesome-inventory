package tenant

import (
	"inv/internal/tenant/domain"
	"inv/internal/tenant/handler"
	"inv/internal/tenant/repository"
	"inv/internal/tenant/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TenantModule struct {
	db            *gorm.DB
	ownerHandler  *handler.TenantOwnerHandler
	tenantHandler *handler.TenantHandler
}

func (m *TenantModule) AutoMigrate() error {
	return m.db.AutoMigrate(
		&domain.Tenant{},
		&domain.TenantOwner{},
	)
}

func NewTenantModule(db *gorm.DB) *TenantModule {
	module := &TenantModule{
		db: db,
	}

	// Run migrations
	if err := module.AutoMigrate(); err != nil {
		panic("failed to migrate tenant tables: " + err.Error())
	}

	// Initialize repositories
	ownerRepo := repository.NewTenantOwnerRepository(db)
	tenantRepo := repository.NewTenantRepository(db)

	// Initialize services
	ownerService := service.NewTenantOwnerService(ownerRepo)
	dbService := service.NewTenantDatabaseService(db)
	tenantService := service.NewTenantService(tenantRepo, dbService)

	// Initialize handlers
	module.ownerHandler = handler.NewTenantOwnerHandler(ownerService)
	module.tenantHandler = handler.NewTenantHandler(tenantService)

	return module
}

func (m *TenantModule) RegisterRoutes(router *gin.Engine) {
	tenantGroup := router.Group("/api/tenants")
	{

		// Owner routes
		tenantGroup.POST("/owners", m.ownerHandler.CreateOwner)
		tenantGroup.PUT("/owners/:id", m.ownerHandler.UpdateOwner)
		tenantGroup.DELETE("/owners/:id", m.ownerHandler.DeleteOwner)
		tenantGroup.GET("/owners/:id", m.ownerHandler.GetOwner)

		// Tenant routes
		tenantGroup.POST("", m.tenantHandler.CreateTenant)
		tenantGroup.PUT("/:id", m.tenantHandler.UpdateTenant)
		tenantGroup.DELETE("/:id", m.tenantHandler.DeleteTenant)
		tenantGroup.GET("/:id", m.tenantHandler.GetTenant)
		tenantGroup.GET("", m.tenantHandler.ListTenants)

	}
}

package app

import (
	"inv/config"
	"inv/internal/tenant"

	// "inv/internal/supplier"
	// "inv/internal/category"
	// "inv/internal/inventory"
	// "inv/internal/product"
	// "inv/internal/user"
	"inv/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	router       *gin.Engine
	db           *gorm.DB
	tenantModule *tenant.TenantModule
	// supplierModule  *supplier.SupplierModule
	// categoryModule  *category.CategoryModule
	// inventoryModule *inventory.InventoryModule
	// productModule   *product.ProductModule
	// userModule   *product.UserModule
}

func New(cfg *config.Config) (*App, error) {
	// Initialize router
	router := gin.Default()

	// Initialize database
	db, err := database.NewConnection(&database.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.Name,
	})
	if err != nil {
		return nil, err
	}

	// Initialize modules
	tenantModule := tenant.NewTenantModule(db)
	// supplierModule := supplier.NewSupplierModule(db)
	// categoryModule := category.NewCategoryModule(db)
	// inventoryModule := inventory.NewInventoryModule(db)
	// productModule := product.NewProductModule(db)
	// userModule := product.NewUserModule(db)

	// Create app instance
	app := &App{
		router:       router,
		db:           db,
		tenantModule: tenantModule,
		// supplierModule:  supplierModule,
		// categoryModule:  categoryModule,
		// inventoryModule: inventoryModule,
		// productModule:   productModule,
		// userModule:      userModule,
	}

	app.setupRoutes()
	return app, nil
}

func (a *App) setupRoutes() {
	// Register routes from all modules
	a.tenantModule.RegisterRoutes(a.router)
	// a.supplierModule.RegisterRoutes(a.router)
	// a.categoryModule.RegisterRoutes(a.router)
	// a.inventoryModule.RegisterRoutes(a.router)
	// a.productModule.RegisterRoutes(a.router)
	// a.userModule.RegisterRoutes(a.router)
}

func (a *App) Run(addr string) error {
	log.Printf("Server is running on %s", addr)
	return a.router.Run(addr)
}

func (a *App) Shutdown() error {
	sqlDB, err := a.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

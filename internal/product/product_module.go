package product

/*

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductModule struct {
	productHandler *handler.ProductHandler
}

func NewProductModule(db *gorm.DB) *ProductModule {
	// Initialize repositories
	productRepo := repository.NewProductRepository(db)

	// Initialize services
	productService := service.NewProductService(productRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productService)

	return &ProductModule{
		productHandler: productHandler,
	}
}

func (m *ProductModule) RegisterRoutes(router *gin.Engine) {
	productGroup := router.Group("/api/products")
	{
		productGroup.POST("", m.productHandler.CreateProduct)
		productGroup.PUT("/:id", m.productHandler.UpdateProduct)
		productGroup.DELETE("/:id", m.productHandler.DeleteProduct)
		productGroup.GET("/:id", m.productHandler.GetProduct)
		productGroup.GET("", m.productHandler.ListProducts)
	}
}

*/

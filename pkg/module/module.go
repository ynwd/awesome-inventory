package module

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module interface {
	RegisterRoutes()
	Migrate() error
}

type ModuleConfig struct {
	Router *gin.Engine
	DB     *gorm.DB
}

package middleware

/*
import (
	"inv/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WithTenantDB(tm database.DatabaseManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "tenant_id not found in context",
			})
			return
		}

		db, err := tm.GetTenantDB(c.Request.Context(), tenantID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get tenant database",
			})
			return
		}

		c.Set("tenant_db", db)
		c.Next()
	}
}
*/

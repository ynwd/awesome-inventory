package middleware

/*
import (
	"strings"

	"github.com/gin-gonic/gin"
)

func TenantResolver() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host
		tenantID := strings.Split(host, ".")[0]

		// Skip for health check endpoint
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		c.Set("tenant_id", tenantID)
		c.Next()
	}
}
*/

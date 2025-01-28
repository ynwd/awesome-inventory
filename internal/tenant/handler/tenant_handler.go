package handler

import (
	"inv/internal/tenant/domain"
	"inv/internal/tenant/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	tenantService service.TenantService
}

func NewTenantHandler(service service.TenantService) *TenantHandler {
	return &TenantHandler{
		tenantService: service,
	}
}

func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var tenant domain.Tenant
	if err := c.ShouldBindJSON(&tenant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values only if fields are empty
	if tenant.DatabaseHost == "" {
		tenant.DatabaseHost = "localhost"
	}
	if tenant.DatabasePort == "" {
		tenant.DatabasePort = "5432"
	}

	if err := h.tenantService.CreateTenant(c.Request.Context(), &tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tenant)
}

func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	id := c.Param("id")
	var tenant domain.Tenant
	if err := c.ShouldBindJSON(&tenant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant.ID = id
	if err := h.tenantService.UpdateTenant(c.Request.Context(), &tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	id := c.Param("id")
	if err := h.tenantService.DeleteTenant(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tenant deleted successfully"})
}

func (h *TenantHandler) GetTenant(c *gin.Context) {
	id := c.Param("id")
	tenant, err := h.tenantService.GetTenant(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) ListTenants(c *gin.Context) {
	tenants, err := h.tenantService.ListTenants(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenants)
}

func (h *TenantHandler) UpdateTenantStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status domain.TenantStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.tenantService.UpdateTenantStatus(c.Request.Context(), id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tenant status updated successfully"})
}

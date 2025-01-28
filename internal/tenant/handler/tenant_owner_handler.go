package handler

import (
	"inv/internal/tenant/domain"
	"inv/internal/tenant/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TenantOwnerHandler struct {
	ownerService service.TenantOwnerService
}

func NewTenantOwnerHandler(service service.TenantOwnerService) *TenantOwnerHandler {
	return &TenantOwnerHandler{
		ownerService: service,
	}
}

func (h *TenantOwnerHandler) CreateOwner(c *gin.Context) {
	var owner domain.TenantOwner
	if err := c.ShouldBindJSON(&owner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ownerService.CreateOwner(c.Request.Context(), &owner); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Do not return the password hash
	owner.Password = ""

	c.JSON(http.StatusCreated, owner)
}

func (h *TenantOwnerHandler) UpdateOwner(c *gin.Context) {
	id := c.Param("id")
	var owner domain.TenantOwner
	if err := c.ShouldBindJSON(&owner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	owner.ID = id
	if err := h.ownerService.UpdateOwner(c.Request.Context(), &owner); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Do not return the password hash
	owner.Password = ""

	c.JSON(http.StatusOK, owner)
}

func (h *TenantOwnerHandler) DeleteOwner(c *gin.Context) {
	id := c.Param("id")
	if err := h.ownerService.DeleteOwner(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "owner deleted successfully"})
}

func (h *TenantOwnerHandler) GetOwner(c *gin.Context) {
	id := c.Param("id")
	owner, err := h.ownerService.GetOwner(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Do not return the password hash
	owner.Password = ""

	c.JSON(http.StatusOK, owner)
}

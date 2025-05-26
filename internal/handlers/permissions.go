package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Anand078/rbac/internal/models"
	"github.com/Anand078/rbac/internal/services"
	"github.com/Anand078/rbac/pkg/utils"
)

type PermissionHandler struct {
	rbacService *services.RBACService
}

func NewPermissionHandler(rbacService *services.RBACService) *PermissionHandler {
	return &PermissionHandler{rbacService: rbacService}
}

func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req models.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	permission, err := h.rbacService.CreatePermission(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Permission created successfully", permission)
}

func (h *PermissionHandler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.rbacService.GetAllPermissions()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permissions retrieved successfully", permissions)
}

func (h *PermissionHandler) GetRolePermissions(c *gin.Context) {
	roleIDStr := c.Param("roleID")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	permissions, err := h.rbacService.GetRolePermissions(roleID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Role permissions retrieved successfully", permissions)
}

func (h *PermissionHandler) GrantPermission(c *gin.Context) {
	var req models.GrantPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	permissionID, err := uuid.Parse(req.PermissionID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	if err := h.rbacService.GrantPermission(roleID, permissionID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permission granted successfully", nil)
}

func (h *PermissionHandler) RevokePermission(c *gin.Context) {
	roleIDStr := c.Param("roleID")
	permissionIDStr := c.Param("permissionID")

	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	if err := h.rbacService.RevokePermission(roleID, permissionID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Permission revoked successfully", nil)
}

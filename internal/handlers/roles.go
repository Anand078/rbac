package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Anand078/rbac/internal/models"
	"github.com/Anand078/rbac/internal/services"
	"github.com/Anand078/rbac/pkg/utils"
)

type RoleHandler struct {
	rbacService *services.RBACService
}

func NewRoleHandler(rbacService *services.RBACService) *RoleHandler {
	return &RoleHandler{rbacService: rbacService}
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req models.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	role, err := h.rbacService.CreateRole(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Role created successfully", role)
}

func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.rbacService.GetAllRoles()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Roles retrieved successfully", roles)
}

func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	roles, err := h.rbacService.GetUserRoles(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User roles retrieved successfully", roles)
}

func (h *RoleHandler) AssignRole(c *gin.Context) {
	var req models.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	if err := h.rbacService.AssignRole(userID, roleID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Role assigned successfully", nil)
}

func (h *RoleHandler) RemoveRole(c *gin.Context) {
	userIDStr := c.Param("userID")
	roleIDStr := c.Param("roleID")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	if err := h.rbacService.RemoveRole(userID, roleID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Role removed successfully", nil)
}

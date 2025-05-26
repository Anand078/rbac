package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/Anand078/rbac/internal/services"
	"github.com/Anand078/rbac/pkg/utils"
)

type AuthMiddleware struct {
	jwtSecret   string
	rbacService *services.RBACService
}

func NewAuthMiddleware(jwtSecret string, rbacService *services.RBACService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:   jwtSecret,
		rbacService: rbacService,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		userID, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID in token")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("email", claims["email"].(string))
		c.Next()
	}
}

func (m *AuthMiddleware) Authorize(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
			c.Abort()
			return
		}

		hasPermission, err := m.rbacService.HasPermission(userID.(uuid.UUID), resource, action)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to check permissions")
			c.Abort()
			return
		}

		if !hasPermission {
			utils.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
			c.Abort()
			return
		}

		roles, err := m.rbacService.GetUserRoles(userID.(uuid.UUID))
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user roles")
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if role.Name == roleName {
				hasRole = true
				break
			}
		}

		if !hasRole {
			utils.ErrorResponse(c, http.StatusForbidden, "Insufficient role")
			c.Abort()
			return
		}

		c.Next()
	}
}

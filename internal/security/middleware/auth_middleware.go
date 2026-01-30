package middleware

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lkgiovani/go-boilerplate/internal/errors"
	"github.com/lkgiovani/go-boilerplate/internal/security/jwt"
)

// AuthMiddleware handles authentication for both web (cookie) and mobile (Authorization header)
type AuthMiddleware struct {
	jwtService *jwt.JwtService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtService *jwt.JwtService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// Authenticate is the middleware function that validates JWT tokens from either cookie or Authorization header
func (m *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	var token string

	// Try to get token from Authorization header first (for mobile)
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			token = parts[1]
		}
	}

	// If no token in header, try to get from cookie (for web)
	if token == "" {
		token = c.Cookies(jwt.AccessTokenCookieName)
	}

	// If still no token found, return unauthorized
	if token == "" {
		return errors.Errorf(errors.EUNAUTHORIZED, "Authentication required")
	}

	// Validate and parse the token
	claims, err := m.jwtService.ParseToken(token)
	if err != nil {
		return errors.Errorf(errors.EUNAUTHORIZED, "Invalid or expired token")
	}

	// Store user information in context for use in handlers
	uid, _ := strconv.ParseInt(claims.ID, 10, 64)
	c.Locals("userID", uid)
	c.Locals("userEmail", claims.Email)
	c.Locals("userName", claims.Name)
	c.Locals("userRoles", claims.Roles)
	c.Locals("userPlan", claims.Plan)
	c.Locals("userAccessMode", claims.AccessMode)

	// Continue to next handler
	return c.Next()
}

// RequireRole creates a middleware that checks if the user has a specific role
func (m *AuthMiddleware) RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roles, ok := c.Locals("userRoles").([]string)
		if !ok {
			return errors.Errorf(errors.EFORBIDDEN, "Access denied")
		}

		// Check if user has the required role
		hasRole := false
		for _, r := range roles {
			if r == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return errors.Errorf(errors.EFORBIDDEN, "Insufficient permissions")
		}

		return c.Next()
	}
}

// RequireAdmin is a convenience middleware for admin-only routes
func (m *AuthMiddleware) RequireAdmin(c *fiber.Ctx) error {
	return m.RequireRole("ADMIN")(c)
}

// RequirePlan creates a middleware that checks if the user has one of the required plans
func (m *AuthMiddleware) RequirePlan(plans ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userPlan, ok := c.Locals("userPlan").(string)
		if !ok {
			return errors.Errorf(errors.EFORBIDDEN, "Access denied")
		}

		// Check if user has one of the required plans
		hasPlan := false
		for _, p := range plans {
			if strings.EqualFold(userPlan, p) {
				hasPlan = true
				break
			}
		}

		if !hasPlan {
			return errors.Errorf(errors.EFORBIDDEN, "Seu plano atual não permite acesso a este recurso")
		}

		return c.Next()
	}
}

// RequireMinPlan ensures the user has at least a certain plan level
func (m *AuthMiddleware) RequireMinPlan(minPlan string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userPlan, ok := c.Locals("userPlan").(string)
		if !ok {
			return errors.Errorf(errors.EFORBIDDEN, "Access denied")
		}

		planLevels := map[string]int{
			"FREE":       1,
			"PRO":        2,
			"ENTERPRISE": 3,
		}

		userLevel := planLevels[strings.ToUpper(userPlan)]
		minLevel := planLevels[strings.ToUpper(minPlan)]

		if userLevel < minLevel {
			return errors.Errorf(errors.EFORBIDDEN, "Este recurso requer um plano %s ou superior", minPlan)
		}

		return c.Next()
	}
}

// RequireWriteAccess ensures the user has a READ_WRITE access mode
func (m *AuthMiddleware) RequireWriteAccess(c *fiber.Ctx) error {
	accessMode, ok := c.Locals("userAccessMode").(string)
	if !ok {
		return errors.Errorf(errors.EFORBIDDEN, "Access denied")
	}

	if !strings.EqualFold(accessMode, "READ_WRITE") {
		return errors.Errorf(errors.EFORBIDDEN, "Sua conta está em modo de apenas leitura ou desativada")
	}

	return c.Next()
}

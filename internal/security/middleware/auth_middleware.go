package middleware

import (
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
		token = c.Cookies(jwt.CookieName)
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
	c.Locals("userID", claims.ID)
	c.Locals("userEmail", claims.Email)
	c.Locals("userName", claims.Name)
	c.Locals("userRoles", claims.Roles)

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

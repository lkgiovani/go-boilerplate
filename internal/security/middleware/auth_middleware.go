package middleware

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lkgiovani/go-boilerplate/internal/errors"
	"github.com/lkgiovani/go-boilerplate/internal/security/jwt"
)

type AuthMiddleware struct {
	jwtService *jwt.JwtService
}

func NewAuthMiddleware(jwtService *jwt.JwtService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	var token string

	authHeader := c.Get("Authorization")
	if authHeader != "" {

		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			token = parts[1]
		}
	}

	if token == "" {
		token = c.Cookies(jwt.AccessTokenCookieName)
	}

	if token == "" {
		return errors.Errorf(errors.EUNAUTHORIZED, "Authentication required")
	}

	claims, err := m.jwtService.ParseToken(token)
	if err != nil {
		return errors.Errorf(errors.EUNAUTHORIZED, "Invalid or expired token")
	}

	uid, _ := strconv.ParseInt(claims.ID, 10, 64)
	c.Locals("userID", uid)
	c.Locals("userEmail", claims.Email)
	c.Locals("userName", claims.Name)
	c.Locals("userRoles", claims.Roles)
	c.Locals("userPlan", claims.Plan)
	c.Locals("userAccessMode", claims.AccessMode)

	return c.Next()
}

func (m *AuthMiddleware) RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roles, ok := c.Locals("userRoles").([]string)
		if !ok {
			return errors.Errorf(errors.EFORBIDDEN, "Access denied")
		}

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

func (m *AuthMiddleware) RequireAdmin(c *fiber.Ctx) error {
	return m.RequireRole("ADMIN")(c)
}

func (m *AuthMiddleware) RequirePlan(plans ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userPlan, ok := c.Locals("userPlan").(string)
		if !ok {
			return errors.Errorf(errors.EFORBIDDEN, "Access denied")
		}

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

package fx

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lkgiovani/go-boilerplate/infra/config"
	"github.com/lkgiovani/go-boilerplate/internal/delivery"
	"github.com/lkgiovani/go-boilerplate/internal/security/middleware"
	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"go.uber.org/fx"
)

var RoutesModule = fx.Module("routes",
	fx.Provide(
		delivery.NewHandler,
		delivery.NewErrorHandler,
		delivery.NewDocsHandler,
		middleware.NewAuthMiddleware,
		delivery.NewMobileAuthHandler,
	),

	fx.Invoke(
		setupRoutes,
	),
)

func setupRoutes(
	lc fx.Lifecycle,
	cfg *config.Config,
	router *fiber.App,

	handler *delivery.Handler,
	docsHandler *delivery.DocsHandler,
	authMiddleware *middleware.AuthMiddleware,

	logger logger.Logger,
) {
	// Health check
	router.Get("/health", handler.HealthCheckHandler)

	// Documentation routes
	docs := router.Group("/docs")
	docs.Get("/", docsHandler.DocsIndex)
	docs.Get("/openapi.yaml", docsHandler.ServeOpenAPI)
	docs.Get("/redoc", docsHandler.ServeRedoc)
	docs.Get("/swagger", docsHandler.ServeSwagger)
	docs.Get("/scalar", docsHandler.ServeScalar)

	// API v1
	v1 := router.Group("/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/refresh", handler.Refresh)
	auth.Post("/logout", handler.Logout)
	auth.Post("/logout-all", authMiddleware.Authenticate, handler.LogoutAll) // Requires authentication
	auth.Post("/signup", handler.Signup)

	// Mobile Auth routes
	mobileAuth := auth.Group("/mobile")
	mobileAuth.Post("/oauth2/google", handler.MobileAuthHandler.AuthenticateWithGoogleMobile)
	mobileAuth.Post("/refresh", handler.MobileAuthHandler.RefreshMobileToken)

	// Public routes
	publicUsers := v1.Group("/users/public")
	publicUsers.Post("/signup", handler.SignupUser)
	publicUsers.Post("/resend-verification", handler.ResendVerification)

	// Upload routes
	uploads := v1.Group("/uploads")
	uploads.Use(authMiddleware.Authenticate)
	uploads.Post("/images", handler.UploadHandler.GetUploadUrl)

	// Email Verification routes
	emailVer := v1.Group("/email-verification")
	emailVer.Post("/send", authMiddleware.Authenticate, handler.EmailVerificationHandler.SendVerificationEmail)
	emailVer.Post("/verify", handler.EmailVerificationHandler.VerifyEmail)
	emailVer.Get("/verify", handler.EmailVerificationHandler.VerifyEmailByQuery)
	emailVer.Post("/resend", handler.EmailVerificationHandler.ResendVerificationEmail)

	// Password Recovery routes
	passRecovery := v1.Group("/password-recovery")
	passRecovery.Post("/request", handler.PasswordRecoveryHandler.RequestPasswordRecovery)
	passRecovery.Post("/verify", handler.PasswordRecoveryHandler.VerifyPasswordRecovery)
	passRecovery.Get("/verify", handler.PasswordRecoveryHandler.VerifyPasswordRecoveryByQuery)
	passRecovery.Post("/reset", handler.PasswordRecoveryHandler.ResetPassword)

	// User routes (require authentication)
	users := v1.Group("/users")
	users.Use(authMiddleware.Authenticate) // Apply authentication middleware to all user routes

	// Authenticated user routes
	users.Get("/me", handler.GetCurrentUser)
	users.Put("/", handler.UpdateUser)
	users.Patch("/password", handler.UpdatePassword)
	users.Patch("/add-image", handler.AddImage)

	// Admin routes (require admin role)
	adminUsers := users.Group("")
	adminUsers.Use(authMiddleware.RequireAdmin) // Apply admin authorization middleware

	adminUsers.Post("/", handler.SaveUser)                            // Create user (admin)
	adminUsers.Get("/", handler.FindAllUsers)                         // List all users (admin)
	adminUsers.Get("/:id", handler.FindUserByID)                      // Get user by ID (admin)
	adminUsers.Get("/email/:email", handler.FindUserByEmail)          // Get user by email (admin)
	adminUsers.Put("/:id", handler.UpdateUserAdmin)                   // Update user (admin)
	adminUsers.Delete("/:id", handler.DeleteUserByID)                 // Delete user (admin)
	adminUsers.Delete("/", handler.DeleteUsersByIDs)                  // Delete multiple users (admin)
	adminUsers.Patch("/:userId/status", handler.ToggleUserStatus)     // Toggle user status (admin)
	adminUsers.Patch("/:id/password", handler.UpdatePasswordAdmin)    // Update user password (admin)
	adminUsers.Patch("/:id/access-mode", handler.UpdateAccessMode)    // Update access mode (admin)
	adminUsers.Patch("/:id/features", handler.UpdateFeatures)         // Update features (admin)
	adminUsers.Patch("/:id/limits", handler.UpdateLimits)             // Update limits (admin)
	adminUsers.Patch("/:id/lifetime-pro", handler.GrantLifetimePro)   // Grant lifetime pro (admin)
	adminUsers.Post("/:id/ensure-metadata", handler.EnsureMetadata)   // Ensure metadata (admin)
	adminUsers.Delete("/:id/lifetime-pro", handler.RevokeLifetimePro) // Revoke lifetime pro (admin)
}

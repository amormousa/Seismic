package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/majoramari/seismic/apps/api/handlers"
)

// Setup registers all API routes on the given Fiber app.
func Setup(app *fiber.App, authHandler *handlers.AuthHandler) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	auth := app.Group("/api/auth")
	auth.Post("/magic-link", authHandler.RequestMagicLink)
	auth.Get("/verify", authHandler.VerifyMagicLink)
	auth.Post("/complete-signup", authHandler.CompleteSignup)
	auth.Post("/refresh", authHandler.RefreshAccessToken)
}

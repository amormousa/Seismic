package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/majoramari/seismic/apps/api/handlers"
	"github.com/majoramari/seismic/apps/api/middleware"
)

func Setup(app *fiber.App, authHandler *handlers.AuthHandler, heartbeatHandler *handlers.HeartbeatHandler, adminHandler *handlers.AdminHandler, statsHandler *handlers.StatsHandler, jwtSecret string, pool *pgxpool.Pool) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	auth := app.Group("/api/auth")
	auth.Post("/magic-link", authHandler.RequestMagicLink)
	auth.Get("/verify", authHandler.VerifyMagicLink)
	auth.Post("/complete-signup", authHandler.CompleteSignup)
	auth.Post("/refresh", authHandler.RefreshAccessToken)

	auth.Get("/apikey", middleware.RequireAuth(jwtSecret), authHandler.GetAPIKey)
	auth.Post("/apikey/regenerate", middleware.RequireAuth(jwtSecret), authHandler.RegenerateAPIKey)
	app.Post("/api/heartbeat", middleware.RequireAPIKey(pool), heartbeatHandler.Receive)
	stats := app.Group("/api/stats", middleware.RequireAuthOrAPIKey(pool, jwtSecret))
	stats.Get("/summary", statsHandler.GetSummary)

	// This is a testing route, not for production use
	app.Post("/api/admin/process-sessions", middleware.RequireAuth(jwtSecret), adminHandler.TriggerSessionProcessing)
}

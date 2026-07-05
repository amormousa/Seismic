package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/majoramari/seismic/apps/api/handlers"
	"github.com/majoramari/seismic/apps/api/middleware"
)

func Setup(app *fiber.App, authHandler *handlers.AuthHandler, heartbeatHandler *handlers.HeartbeatHandler, adminHandler *handlers.AdminHandler, statsHandler *handlers.StatsHandler, jwtSecret string, pool *pgxpool.Pool) {
	app.Get("/health", handlers.HealthCheck)
	auth := app.Group("/api/auth")
	auth.Get("/verify", authHandler.VerifyMagicLink)
	auth.Post("/complete-signup", authHandler.CompleteSignup)
	auth.Post("/refresh", authHandler.RefreshAccessToken)

	auth.Get("/apikey", middleware.RequireAuth(jwtSecret), authHandler.GetAPIKey)
	auth.Post("/apikey/regenerate", middleware.RequireAuth(jwtSecret), authHandler.RegenerateAPIKey)
	auth.Post("/magic-link", middleware.AuthRateLimit(), authHandler.RequestMagicLink)
	app.Post("/api/heartbeat", middleware.HeartbeatRateLimit(), middleware.RequireAPIKey(pool), heartbeatHandler.Receive)
	auth.Get("/check-username", authHandler.CheckUsername)

	stats := app.Group("/api/stats", middleware.RequireAuthOrAPIKey(pool, jwtSecret))
	stats.Get("/summary", statsHandler.GetSummary)
	stats.Get("/languages", statsHandler.GetLanguages)
	stats.Get("/heatmap", statsHandler.GetHeatmap)

	leaderboardHandler := &handlers.LeaderboardHandler{Pool: pool}
	app.Get("/api/leaderboard", leaderboardHandler.GetLeaderboard)

	// This is a testing route, not for production use
	app.Post("/api/admin/process-sessions", middleware.RequireAuth(jwtSecret), adminHandler.TriggerSessionProcessing)
}

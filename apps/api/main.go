// main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/majoramari/seismic/apps/api/config"
	"github.com/majoramari/seismic/apps/api/db"
	"github.com/majoramari/seismic/apps/api/handlers"
	"github.com/majoramari/seismic/apps/api/routes"
	"github.com/majoramari/seismic/apps/api/services"
)

func main() {
	cfg := config.Load()
	pool := db.Connect(cfg.DatabaseURL)
	defer pool.Close()

	if err := db.RunMigrations(pool); err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}

	app := fiber.New()

	authHandler := &handlers.AuthHandler{
		Pool:      pool,
		JWTSecret: cfg.JWTSecret,
		EmailCfg: services.EmailConfig{
			Host:     cfg.SMTPHost,
			Port:     cfg.SMTPPort,
			Username: cfg.SMTPUser,
			Password: cfg.SMTPPass,
			AppURL:   cfg.AppURL,
		},
	}

	heartbeatHandler := &handlers.HeartbeatHandler{Pool: pool}

	routes.Setup(app, authHandler, heartbeatHandler, cfg.JWTSecret, pool)

	log.Printf("Seismic API starting on port %s\n", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

package handlers

import "github.com/gofiber/fiber/v2"

// HealthCheck godoc
// @Summary      Health check
// @Description  Returns ok if the server is running
// @Tags         system
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /health [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

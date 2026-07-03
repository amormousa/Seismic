package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/majoramari/seismic/apps/api/helpers"
	"github.com/majoramari/seismic/apps/api/services"
)

type AdminHandler struct {
	Pool *pgxpool.Pool
}

// TriggerSessionProcessing Manually triggers the session processor, for testing only.
func (h *AdminHandler) TriggerSessionProcessing(c *fiber.Ctx) error {
	ctx := c.Context()
	if err := services.ProcessSessions(ctx, h.Pool); err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to process sessions")
	}
	return helpers.Success(c, "Sessions processed", nil)
}

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

// TriggerSessionProcessing godoc
// @Summary      Manually trigger session processing
// @Description  Runs the session processor immediately instead of waiting for the 5 minute timer. For testing only.
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} helpers.APIResponse
// @Failure      401 {object} helpers.APIResponse
// @Router       /api/admin/process-sessions [post]
func (h *AdminHandler) TriggerSessionProcessing(c *fiber.Ctx) error {
	ctx := c.Context()
	if err := services.ProcessSessions(ctx, h.Pool); err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to process sessions")
	}
	return helpers.Success(c, "Sessions processed", nil)
}

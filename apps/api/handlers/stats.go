package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/majoramari/seismic/apps/api/helpers"
	"github.com/majoramari/seismic/apps/api/models"
)

type StatsHandler struct {
	Pool *pgxpool.Pool
}

// GetSummary handles GET /api/stats/summary?range=today|week|month|all
func (h *StatsHandler) GetSummary(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	rangeParam := c.Query("range", "today")

	ctx := c.Context()
	summary, err := models.GetStatsSummary(ctx, h.Pool, userID, models.RangeSQL(rangeParam))
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to fetch stats")
	}

	return helpers.Success(c, "Stats retrieved", summary)
}

func (h *StatsHandler) GetLanguages(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	rangeParam := c.Query("range", "week")

	ctx := c.Context()
	stats, err := models.GetLanguageBreakdown(ctx, h.Pool, userID, models.RangeSQL(rangeParam))
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to fetch stats")
	}

	return helpers.Success(c, "Language breakdown retrieved", stats)
}

func (h *StatsHandler) GetHeatmap(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	ctx := c.Context()
	days, err := models.GetHeatmap(ctx, h.Pool, userID)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to fetch heatmap")
	}

	return helpers.Success(c, "Heatmap retrieved", days)
}

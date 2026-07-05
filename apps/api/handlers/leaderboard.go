package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/majoramari/seismic/apps/api/helpers"
	"github.com/majoramari/seismic/apps/api/models"
)

type LeaderboardHandler struct {
	Pool *pgxpool.Pool
}

// GetLeaderboard godoc
// @Summary      Get leaderboard
// @Description  Returns ranked users by total coding time within a range. If authenticated, includes the viewer's own rank.
// @Tags         leaderboard
// @Produce      json
// @Param        range query string false "today, week, month, or all" default(today)
// @Success      200 {object} helpers.APIResponse
// @Router       /api/leaderboard [get]
func (h *LeaderboardHandler) GetLeaderboard(c *fiber.Ctx) error {
	rangeParam := c.Query("range", "today")

	currentUserID, _ := c.Locals("userID").(string)

	ctx := c.Context()
	result, err := models.GetLeaderboard(ctx, h.Pool, models.RangeSQL(rangeParam), 50, currentUserID)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to fetch leaderboard")
	}

	return helpers.Success(c, "Leaderboard retrieved", result)
}

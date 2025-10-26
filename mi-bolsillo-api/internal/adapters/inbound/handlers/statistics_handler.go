package handlers

import (
	"net/http"
	"strconv"

	"github.com/KKogaa/mi-bolsillo-api/internal/core/services"
	"github.com/labstack/echo/v4"
)

type StatisticsHandler struct {
	statisticsService *services.StatisticsService
	accountLinkService *services.AccountLinkService
}

func NewStatisticsHandler(statisticsService *services.StatisticsService, accountLinkService *services.AccountLinkService) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
		accountLinkService: accountLinkService,
	}
}

// GetDashboardStatistics godoc
// @Summary Get dashboard statistics
// @Description Returns comprehensive statistics including monthly, weekly, and category breakdowns
// @Tags statistics
// @Produce json
// @Param months query int false "Number of months to include" default(6)
// @Success 200 {object} dtos.DashboardStatistics
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /statistics/dashboard [get]
func (h *StatisticsHandler) GetDashboardStatistics(c echo.Context) error {
	// Get user ID from JWT (set by Clerk middleware)
	clerkID, ok := c.Get("userID").(string)
	if !ok || clerkID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized - invalid user ID",
		})
	}

	// Get or create user by Clerk ID
	user, err := h.accountLinkService.GetOrCreateUserByClerkID(clerkID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user information",
		})
	}

	// Get months parameter (default to 6)
	months := 6
	if monthsParam := c.QueryParam("months"); monthsParam != "" {
		if m, err := strconv.Atoi(monthsParam); err == nil && m > 0 && m <= 24 {
			months = m
		}
	}

	// Get statistics
	stats, err := h.statisticsService.GetDashboardStatistics(user.UserID, months)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch statistics",
		})
	}

	return c.JSON(http.StatusOK, stats)
}

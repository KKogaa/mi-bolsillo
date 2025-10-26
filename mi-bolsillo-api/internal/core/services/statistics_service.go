package services

import (
	"fmt"
	"time"

	"github.com/KKogaa/mi-bolsillo-api/internal/core/entities"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/ports"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/services/dtos"
)

type StatisticsService struct {
	billRepo ports.BillRepository
}

func NewStatisticsService(billRepo ports.BillRepository) *StatisticsService {
	return &StatisticsService{
		billRepo: billRepo,
	}
}

// GetDashboardStatistics returns comprehensive statistics for the user's dashboard
func (s *StatisticsService) GetDashboardStatistics(userID string, months int) (*dtos.DashboardStatistics, error) {
	// Get all bills for the user
	bills, err := s.billRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bills: %w", err)
	}

	// Calculate statistics
	monthlyStats := s.calculateMonthlyStatistics(bills, months)
	weeklyStats := s.calculateWeeklyStatistics(bills, 8) // Last 8 weeks
	categoryStats := s.calculateCategoryStatistics(bills)

	// Calculate totals
	var totalPEN, totalUSD float64
	for _, bill := range bills {
		totalPEN += bill.AmountPen
		totalUSD += bill.AmountUsd
	}

	return &dtos.DashboardStatistics{
		MonthlyStats:  monthlyStats,
		WeeklyStats:   weeklyStats,
		CategoryStats: categoryStats,
		TotalPEN:      totalPEN,
		TotalUSD:      totalUSD,
		TotalBills:    len(bills),
	}, nil
}

// calculateMonthlyStatistics calculates monthly spending statistics
func (s *StatisticsService) calculateMonthlyStatistics(bills []*entities.Bill, months int) []dtos.MonthlyStatistics {
	monthlyMap := make(map[string]*dtos.MonthlyStatistics)

	now := time.Now()
	// Initialize last N months
	for i := 0; i < months; i++ {
		targetDate := now.AddDate(0, -i, 0)
		monthKey := targetDate.Format("2006-01")
		monthlyMap[monthKey] = &dtos.MonthlyStatistics{
			Month:     monthKey,
			Year:      targetDate.Year(),
			MonthNum:  int(targetDate.Month()),
			TotalPEN:  0,
			TotalUSD:  0,
			BillCount: 0,
		}
	}

	// Aggregate bills by month
	for _, bill := range bills {
		monthKey := bill.Date.Format("2006-01")
		if stats, exists := monthlyMap[monthKey]; exists {
			stats.TotalPEN += bill.AmountPen
			stats.TotalUSD += bill.AmountUsd
			stats.BillCount++
		}
	}

	// Convert map to slice and sort by date (newest first)
	result := make([]dtos.MonthlyStatistics, 0, len(monthlyMap))
	for i := 0; i < months; i++ {
		targetDate := now.AddDate(0, -i, 0)
		monthKey := targetDate.Format("2006-01")
		if stats, exists := monthlyMap[monthKey]; exists {
			result = append(result, *stats)
		}
	}

	return result
}

// calculateWeeklyStatistics calculates weekly spending statistics
func (s *StatisticsService) calculateWeeklyStatistics(bills []*entities.Bill, weeks int) []dtos.WeeklyStatistics {
	weeklyMap := make(map[string]*dtos.WeeklyStatistics)

	now := time.Now()

	// Initialize last N weeks
	for i := 0; i < weeks; i++ {
		weekStart := getWeekStart(now.AddDate(0, 0, -7*i))
		weekEnd := weekStart.AddDate(0, 0, 6)
		weekKey := weekStart.Format("2006-01-02")

		weeklyMap[weekKey] = &dtos.WeeklyStatistics{
			WeekStart: weekStart,
			WeekEnd:   weekEnd,
			WeekLabel: fmt.Sprintf("%s - %s", weekStart.Format("Jan 2"), weekEnd.Format("Jan 2")),
			TotalPEN:  0,
			TotalUSD:  0,
			BillCount: 0,
		}
	}

	// Aggregate bills by week
	for _, bill := range bills {
		weekStart := getWeekStart(bill.Date)
		weekKey := weekStart.Format("2006-01-02")

		if stats, exists := weeklyMap[weekKey]; exists {
			stats.TotalPEN += bill.AmountPen
			stats.TotalUSD += bill.AmountUsd
			stats.BillCount++
		}
	}

	// Convert map to slice and sort by date (newest first)
	result := make([]dtos.WeeklyStatistics, 0, len(weeklyMap))
	for i := 0; i < weeks; i++ {
		weekStart := getWeekStart(now.AddDate(0, 0, -7*i))
		weekKey := weekStart.Format("2006-01-02")
		if stats, exists := weeklyMap[weekKey]; exists {
			result = append(result, *stats)
		}
	}

	return result
}

// calculateCategoryStatistics calculates spending by category
func (s *StatisticsService) calculateCategoryStatistics(bills []*entities.Bill) []dtos.CategoryStatistics {
	categoryMap := make(map[string]*dtos.CategoryStatistics)
	var totalPEN float64

	// Aggregate by category
	for _, bill := range bills {
		category := bill.Category
		if category == "" {
			category = "Uncategorized"
		}

		if _, exists := categoryMap[category]; !exists {
			categoryMap[category] = &dtos.CategoryStatistics{
				Category:  category,
				TotalPEN:  0,
				TotalUSD:  0,
				BillCount: 0,
			}
		}

		categoryMap[category].TotalPEN += bill.AmountPen
		categoryMap[category].TotalUSD += bill.AmountUsd
		categoryMap[category].BillCount++
		totalPEN += bill.AmountPen
	}

	// Calculate percentages and convert to slice
	result := make([]dtos.CategoryStatistics, 0, len(categoryMap))
	for _, stats := range categoryMap {
		if totalPEN > 0 {
			stats.Percentage = (stats.TotalPEN / totalPEN) * 100
		}
		result = append(result, *stats)
	}

	// Sort by total (highest first)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].TotalPEN > result[i].TotalPEN {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

// getWeekStart returns the Monday of the week for a given date
func getWeekStart(date time.Time) time.Time {
	weekday := date.Weekday()
	// Go's Sunday = 0, Monday = 1, ..., Saturday = 6
	// We want Monday as start of week
	daysFromMonday := int(weekday) - 1
	if weekday == time.Sunday {
		daysFromMonday = 6
	}

	weekStart := date.AddDate(0, 0, -daysFromMonday)
	return time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
}

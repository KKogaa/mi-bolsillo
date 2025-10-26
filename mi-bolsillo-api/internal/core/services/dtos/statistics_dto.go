package dtos

import "time"

// MonthlyStatistics represents spending statistics for a month
type MonthlyStatistics struct {
	Month      string  `json:"month"`      // Format: "2024-01"
	TotalPEN   float64 `json:"totalPen"`   // Total spent in PEN
	TotalUSD   float64 `json:"totalUsd"`   // Total spent in USD
	BillCount  int     `json:"billCount"`  // Number of bills
	Year       int     `json:"year"`       // Year
	MonthNum   int     `json:"monthNum"`   // Month number (1-12)
}

// WeeklyStatistics represents spending statistics for a week
type WeeklyStatistics struct {
	WeekStart  time.Time `json:"weekStart"`  // Start of the week (Monday)
	WeekEnd    time.Time `json:"weekEnd"`    // End of the week (Sunday)
	WeekLabel  string    `json:"weekLabel"`  // Format: "Week 1 (Jan 1 - Jan 7)"
	TotalPEN   float64   `json:"totalPen"`   // Total spent in PEN
	TotalUSD   float64   `json:"totalUsd"`   // Total spent in USD
	BillCount  int       `json:"billCount"`  // Number of bills
}

// CategoryStatistics represents spending statistics by category
type CategoryStatistics struct {
	Category  string  `json:"category"`   // Category name
	TotalPEN  float64 `json:"totalPen"`   // Total spent in PEN
	TotalUSD  float64 `json:"totalUsd"`   // Total spent in USD
	BillCount int     `json:"billCount"`  // Number of bills
	Percentage float64 `json:"percentage"` // Percentage of total spending
}

// DashboardStatistics represents overall dashboard statistics
type DashboardStatistics struct {
	MonthlyStats  []MonthlyStatistics  `json:"monthlyStats"`
	WeeklyStats   []WeeklyStatistics   `json:"weeklyStats"`
	CategoryStats []CategoryStatistics `json:"categoryStats"`
	TotalPEN      float64              `json:"totalPen"`
	TotalUSD      float64              `json:"totalUsd"`
	TotalBills    int                  `json:"totalBills"`
}

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/KKogaa/mi-bolsillo-api/config"
	"github.com/KKogaa/mi-bolsillo-api/internal/adapters/inbound/handlers"
	custommiddleware "github.com/KKogaa/mi-bolsillo-api/internal/adapters/inbound/middleware"
	"github.com/KKogaa/mi-bolsillo-api/internal/adapters/outbound/grok"
	"github.com/KKogaa/mi-bolsillo-api/internal/adapters/outbound/repositories"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/services"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	_ "github.com/KKogaa/mi-bolsillo-api/docs" // Import generated docs
)

// @title Mi Bolsillo API
// @version 1.0
// @description API for managing bills and expenses with multi-currency support
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@mibolsillo.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token from Clerk authentication

func connectToDatabase(cfg *config.Config) (*sqlx.DB, error) {
	dbUrl := fmt.Sprintf("%s?authToken=%s", cfg.DatabaseUrl, cfg.DatabaseToken)
	db, err := sqlx.Connect("libsql", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to libsql database")
	return db, nil
}

func runMigrations(db *sqlx.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id TEXT PRIMARY KEY,
			clerk_id TEXT UNIQUE,
			telegram_id INTEGER UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create account_link_otps table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS account_link_otps (
			otp_code TEXT PRIMARY KEY,
			telegram_id INTEGER NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create account_link_otps table: %w", err)
	}

	// Create bills table with source field
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bills (
			bill_id TEXT PRIMARY KEY,
			amount_pen REAL NOT NULL,
			amount_usd REAL NOT NULL,
			description TEXT,
			category TEXT,
			currency TEXT NOT NULL,
			user_id TEXT NOT NULL,
			source TEXT NOT NULL DEFAULT 'web',
			date DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create bills table: %w", err)
	}

	// Create expenses table with source field
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			expense_id TEXT PRIMARY KEY,
			amount_pen REAL NOT NULL,
			amount_usd REAL NOT NULL,
			exchange_rate REAL NOT NULL,
			currency TEXT NOT NULL,
			description TEXT,
			category TEXT,
			date TEXT NOT NULL,
			bill_id TEXT,
			user_id TEXT NOT NULL,
			source TEXT NOT NULL DEFAULT 'web',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (bill_id) REFERENCES bills(bill_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create expenses table: %w", err)
	}

	// Add source column to existing bills table if it doesn't exist
	// SQLite doesn't have a simple way to check if column exists, so we try to add it
	// and ignore errors if it already exists
	_, _ = db.Exec(`ALTER TABLE bills ADD COLUMN source TEXT NOT NULL DEFAULT 'web'`)

	// Add source column to existing expenses table if it doesn't exist
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN source TEXT NOT NULL DEFAULT 'web'`)

	log.Println("Successfully ran database migrations")
	return nil
}

func migrateExistingData(db *sqlx.DB) error {
	// This function migrates existing bills and creates user records for them
	// It's idempotent and can be run multiple times safely

	log.Println("Starting data migration for existing bills...")

	// Find all distinct user_ids in bills table
	var userIDs []string
	err := db.Select(&userIDs, `SELECT DISTINCT user_id FROM bills`)
	if err != nil {
		log.Printf("Failed to get user IDs from bills: %v", err)
		// Don't fail if there are no bills yet
		return nil
	}

	for _, userID := range userIDs {
		// Check if this user already exists
		var count int
		err := db.Get(&count, `SELECT COUNT(*) FROM users WHERE user_id = ?`, userID)
		if err != nil {
			log.Printf("Failed to check user existence: %v", err)
			continue
		}

		if count > 0 {
			// User already exists, skip
			continue
		}

		// Determine if this is a telegram or clerk user based on user_id format
		// Telegram users have format "tg_{telegramID}"
		// Clerk users have other formats
		if len(userID) > 3 && userID[:3] == "tg_" {
			// This is a telegram user (old format)
			// Parse telegram ID from the user_id
			var telegramID int64
			_, err := fmt.Sscanf(userID, "tg_%d", &telegramID)
			if err != nil {
				log.Printf("Failed to parse telegram ID from %s: %v", userID, err)
				continue
			}

			// Create user with telegram ID
			_, err = db.Exec(`
				INSERT INTO users (user_id, telegram_id, created_at, updated_at)
				VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			`, userID, telegramID)
			if err != nil {
				log.Printf("Failed to create telegram user %s: %v", userID, err)
				continue
			}
			log.Printf("Created user record for telegram user: %s", userID)
		} else {
			// This is a clerk user
			_, err = db.Exec(`
				INSERT INTO users (user_id, clerk_id, created_at, updated_at)
				VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			`, userID, userID)
			if err != nil {
				log.Printf("Failed to create clerk user %s: %v", userID, err)
				continue
			}
			log.Printf("Created user record for clerk user: %s", userID)
		}
	}

	log.Println("Data migration completed")
	return nil
}

func main() {
	cfg := config.LoadConfig()
	db, err := connectToDatabase(cfg)
	if err != nil {
		log.Fatal("Database connection error", "error", err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		log.Fatal("Database migration error", "error", err)
	}

	// Migrate existing data
	if err := migrateExistingData(db); err != nil {
		log.Printf("Warning: Data migration encountered errors: %v", err)
		// Don't fail startup on data migration errors
	}

	// Initialize repositories
	billRepo := repositories.NewBillRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)
	userRepo := repositories.NewUserRepository(db)
	otpRepo := repositories.NewOTPRepository(db)

	// Initialize services
	billWithExpensesService := services.NewBillWithExpensesService(billRepo, expenseRepo)
	accountLinkService := services.NewAccountLinkService(userRepo, otpRepo, billRepo, expenseRepo, cfg.OTPExpirationMinutes)
	statisticsService := services.NewStatisticsService(billRepo)

	// Initialize Grok client
	grokClient := grok.NewGrokClient(cfg.GrokAPIKey)

	// Initialize handlers
	billWithExpensesHandler := handlers.NewBillWithExpensesHandler(billWithExpensesService, accountLinkService)
	billUploadHandler := handlers.NewBillUploadHandler(grokClient, billWithExpensesService, accountLinkService)
	authHandler := handlers.NewAuthHandler(accountLinkService)
	statisticsHandler := handlers.NewStatisticsHandler(statisticsService, accountLinkService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	// Swagger documentation route (public)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Protected routes group with Clerk authentication
	api := e.Group("")
	api.Use(custommiddleware.ClerkAuthWithConfig(cfg.ClerkJWKSUrl))

	// Register routes
	api.POST("/bills", billWithExpensesHandler.CreateBillWithExpenses)
	api.POST("/bills/upload", billUploadHandler.UploadBillPhoto)
	api.GET("/bills", billWithExpensesHandler.ListBills)
	api.GET("/bills/:id", billWithExpensesHandler.GetBillByID)
	api.DELETE("/bills/:id", billWithExpensesHandler.DeleteBillByID)
	api.POST("/auth/verify-otp", authHandler.VerifyOTP)
	api.GET("/auth/link-status", authHandler.GetLinkStatus)
	api.GET("/statistics/dashboard", statisticsHandler.GetDashboardStatistics)

	// Use PORT from config (Render will set this automatically)
	port := cfg.Port
	if port == "" {
		port = "8080" // Default fallback
	}
	serverAddr := fmt.Sprintf(":%s", port)

	log.Printf("Starting server on %s", serverAddr)
	if err := e.Start(serverAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("failed to start server", "error", err)
	}

}

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/KKogaa/mi-bolsillo-api/config"
	"github.com/KKogaa/mi-bolsillo-api/internal/adapters/inbound/handlers"
	custommiddleware "github.com/KKogaa/mi-bolsillo-api/internal/adapters/inbound/middleware"
	"github.com/KKogaa/mi-bolsillo-api/internal/adapters/outbound/repositories"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/services"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

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
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS bills (
			bill_id TEXT PRIMARY KEY,
			amount_pen REAL NOT NULL,
			amount_usd REAL NOT NULL,
			description TEXT,
			category TEXT,
			currency TEXT NOT NULL,
			user_id TEXT NOT NULL,
			date DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create bills table: %w", err)
	}

	// Create expenses table based on entities.Expense
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
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (bill_id) REFERENCES bills(bill_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create expenses table: %w", err)
	}

	log.Println("Successfully ran database migrations")
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

	// Initialize repositories
	billRepo := repositories.NewBillRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)

	// Initialize services
	billWithExpensesService := services.NewBillWithExpensesService(billRepo, expenseRepo)

	// Initialize handlers
	billWithExpensesHandler := handlers.NewBillWithExpensesHandler(billWithExpensesService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Protected routes group with Clerk authentication
	api := e.Group("")
	api.Use(custommiddleware.ClerkAuthWithConfig(cfg.ClerkJWKSUrl))

	// Register routes
	api.POST("/bills", billWithExpensesHandler.CreateBillWithExpenses)
	api.GET("/bills", billWithExpensesHandler.ListBills)
	api.GET("/bills/:id", billWithExpensesHandler.GetBillByID)
	api.DELETE("/bills/:id", billWithExpensesHandler.DeleteBillByID)

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("failed to start server", "error", err)
	}

}

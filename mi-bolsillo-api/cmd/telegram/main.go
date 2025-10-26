package main

import (
	"fmt"
	"log"
	"time"

	"github.com/KKogaa/mi-bolsillo-api/config"
	"github.com/KKogaa/mi-bolsillo-api/internal/adapters/inbound/handlers/telegram"
	"github.com/KKogaa/mi-bolsillo-api/internal/adapters/outbound/grok"
	"github.com/KKogaa/mi-bolsillo-api/internal/adapters/outbound/repositories"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/services"
	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	tele "gopkg.in/telebot.v3"
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

func main() {
	cfg := config.LoadConfig()

	// Load messages
	messages, err := telegram.LoadMessages("config/messages.json")
	if err != nil {
		log.Fatal("Failed to load messages:", err)
	}

	// Connect to database
	db, err := connectToDatabase(cfg)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	// Initialize repositories
	billRepo := repositories.NewBillRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)
	userRepo := repositories.NewUserRepository(db)
	otpRepo := repositories.NewOTPRepository(db)

	// Initialize services
	billWithExpensesService := services.NewBillWithExpensesService(billRepo, expenseRepo)
	accountLinkService := services.NewAccountLinkService(userRepo, otpRepo, billRepo, expenseRepo, cfg.OTPExpirationMinutes)

	// Initialize Grok client (implements IntentDetector interface)
	grokClient := grok.NewGrokClient(cfg.GrokAPIKey)

	// Create bot handler
	botHandler := telegram.NewBotHandler(
		grokClient, // GrokClient implements ports.IntentDetector
		billWithExpensesService,
		accountLinkService,
		grokClient,
		messages,
	)

	// Create bot
	pref := tele.Settings{
		Token:  cfg.TelegramBotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	// Register handlers
	bot.Handle("/start", botHandler.HandleStart)
	bot.Handle("/link", botHandler.HandleLink)
	bot.Handle(tele.OnText, botHandler.HandleText)
	bot.Handle(tele.OnPhoto, botHandler.HandlePhoto)

	log.Println("Telegram bot started successfully using long polling")
	bot.Start()
}

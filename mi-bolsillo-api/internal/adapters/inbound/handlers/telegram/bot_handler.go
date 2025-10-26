package telegram

import (
	"fmt"
	"io"
	"log"
	"time"

	handlerdtos "github.com/KKogaa/mi-bolsillo-api/internal/adapters/inbound/handlers/dtos"
	"github.com/KKogaa/mi-bolsillo-api/internal/adapters/inbound/handlers/mappers"
	"github.com/KKogaa/mi-bolsillo-api/internal/adapters/outbound/grok"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/domain/entities"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/ports"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/services"
	servicedtos "github.com/KKogaa/mi-bolsillo-api/internal/core/services/dtos"
	tele "gopkg.in/telebot.v3"
)

// BotHandler holds all dependencies for the bot handlers
type BotHandler struct {
	intentDetector          ports.IntentDetector
	billWithExpensesService *services.BillWithExpensesService
	accountLinkService      *services.AccountLinkService
	grokClient              *grok.GrokClient
	messages                *Messages
}

// NewBotHandler creates a new BotHandler instance
func NewBotHandler(
	intentDetector ports.IntentDetector,
	billWithExpensesService *services.BillWithExpensesService,
	accountLinkService *services.AccountLinkService,
	grokClient *grok.GrokClient,
	messages *Messages,
) *BotHandler {
	return &BotHandler{
		intentDetector:          intentDetector,
		billWithExpensesService: billWithExpensesService,
		accountLinkService:      accountLinkService,
		grokClient:              grokClient,
		messages:                messages,
	}
}

func (h *BotHandler) HandleStart(c tele.Context) error {
	return c.Send(h.messages.Welcome, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

func (h *BotHandler) HandleLink(c tele.Context) error {
	telegramID := c.Sender().ID

	// Get or create user by Telegram ID to ensure user exists before linking
	_, err := h.accountLinkService.GetOrCreateUserByTelegramID(telegramID)
	if err != nil {
		log.Printf("Failed to get or create user for Telegram ID %d: %v", telegramID, err)
		return c.Send(h.messages.LinkAccountError, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	}

	// Generate OTP
	otpCode, err := h.accountLinkService.GenerateOTP(telegramID)
	if err != nil {
		log.Printf("Failed to generate OTP for user %d: %v", telegramID, err)
		return c.Send(h.messages.LinkAccountError, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	}

	// Send OTP to user
	message := fmt.Sprintf(h.messages.LinkAccountOTP, otpCode)
	return c.Send(message, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

func (h *BotHandler) HandleText(c tele.Context) error {
	telegramID := c.Sender().ID
	text := c.Text()

	// Get or create user by Telegram ID
	user, err := h.accountLinkService.GetOrCreateUserByTelegramID(telegramID)
	if err != nil {
		log.Printf("Failed to get user for Telegram ID %d: %v", telegramID, err)
		return c.Send(h.messages.ErrorProcessingMsg)
	}

	log.Printf("Received text from user %d: %s", telegramID, text)

	// Detect intent
	intent, err := h.intentDetector.DetectIntent(text)
	if err != nil {
		log.Printf("Failed to detect intent: %v", err)
		return c.Send(h.messages.ErrorUnderstand)
	}

	log.Printf("Detected intent: %s (confidence: %.2f)", intent.Type, intent.Confidence)

	switch intent.Type {
	case entities.IntentListBills:
		return h.handleListBills(c, user.UserID, intent)
	case entities.IntentSummaryBills:
		return h.handleSummaryBills(c, user.UserID, intent)
	case entities.IntentCreateExpense:
		return h.handleCreateExpense(c, user.UserID, intent)
	case entities.IntentUnknown:
		fallthrough
	default:
		return c.Send(h.messages.UnknownIntent, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	}
}

func (h *BotHandler) HandlePhoto(c tele.Context) error {
	telegramID := c.Sender().ID

	// Get or create user by Telegram ID
	user, err := h.accountLinkService.GetOrCreateUserByTelegramID(telegramID)
	if err != nil {
		log.Printf("Failed to get user for Telegram ID %d: %v", telegramID, err)
		return c.Send(h.messages.ErrorProcessingMsg)
	}

	log.Printf("Received photo from user %d", telegramID)

	// Send processing message
	if err := c.Send(h.messages.ProcessingImage); err != nil {
		log.Printf("Failed to send processing message: %v", err)
	}

	// Get the photo
	photo := c.Message().Photo
	if photo == nil {
		return c.Send(h.messages.ErrorRetrieveImage)
	}

	// Download the file
	file, err := c.Bot().FileByID(photo.FileID)
	if err != nil {
		log.Printf("Failed to get file: %v", err)
		return c.Send(h.messages.ErrorRetrieveImage)
	}

	reader, err := c.Bot().File(&file)
	if err != nil {
		log.Printf("Failed to download file: %v", err)
		return c.Send(h.messages.ErrorDownloadImage)
	}
	defer reader.Close()

	// Read the image data into a byte slice
	imageData, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("Failed to read image data: %v", err)
		return c.Send(h.messages.ErrorReadImage)
	}

	// Parse the bill using Grok
	parsedData, err := h.grokClient.ParseBillImage(imageData)
	if err != nil {
		log.Printf("Failed to parse bill image: %v", err)
		return c.Send(h.messages.ErrorParseBill)
	}

	// Parse date string to time.Time
	billDate, err := time.Parse("2006-01-02", parsedData.Date)
	if err != nil {
		billDate = time.Now()
	}

	// Create bill with expenses
	// Set exchange rate based on currency
	exchangeRate := 3.75 // Default PEN to USD exchange rate
	if parsedData.Currency == "USD" {
		exchangeRate = 3.75 // USD to PEN
	}

	handlerDTO := handlerdtos.CreateBillWithExpensesRequest{
		UserID:       user.UserID,
		Source:       "telegram",
		Description:  parsedData.MerchantName,
		Category:     "General",
		Currency:     parsedData.Currency,
		ExchangeRate: exchangeRate,
		Date:         billDate,
		Expenses:     make([]handlerdtos.CreateExpenseForBill, len(parsedData.Items)),
	}

	for i, item := range parsedData.Items {
		handlerDTO.Expenses[i] = handlerdtos.CreateExpenseForBill{
			Description: item.Description,
			Amount:      item.Amount,
			Category:    item.Category,
			Date:        parsedData.Date,
		}
	}

	serviceDTO := mappers.ToCreateBillWithExpensesServiceDTO(handlerDTO)
	bill, _, err := h.billWithExpensesService.CreateBillWithExpenses(serviceDTO)
	if err != nil {
		log.Printf("Failed to create bill: %v", err)
		return c.Send(h.messages.ErrorSaveBill)
	}

	// Send success message
	responseMsg := fmt.Sprintf(h.messages.BillSaved,
		parsedData.MerchantName,
		parsedData.Currency,
		parsedData.TotalAmount,
		parsedData.Date,
		len(parsedData.Items),
	)

	log.Printf("Bill created successfully: %s", bill.BillId)
	return c.Send(responseMsg, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

func (h *BotHandler) handleListBills(c tele.Context, userID string, intent *entities.Intent) error {
	bills, err := h.billWithExpensesService.ListBillsByUserID(userID)
	if err != nil {
		log.Printf("Failed to list bills: %v", err)
		return c.Send(h.messages.ErrorRetrieveBills)
	}

	if len(bills) == 0 {
		return c.Send(h.messages.NoBills)
	}

	// Determine limit
	limit := 5
	if limitParam, ok := intent.Parameters["limit"].(float64); ok {
		limit = int(limitParam)
	}

	if limit > len(bills) {
		limit = len(bills)
	}

	// Build response message
	responseMsg := fmt.Sprintf(h.messages.BillsListHeader, limit, len(bills))

	for i := 0; i < limit; i++ {
		bill := bills[i]
		responseMsg += fmt.Sprintf("*%d.* %s\n", i+1, bill.Description)
		responseMsg += fmt.Sprintf("   💰 %s %.2f (PEN %.2f / USD %.2f)\n", bill.Currency,
			getAmountInCurrency(bill, bill.Currency), bill.AmountPen, bill.AmountUsd)
		responseMsg += fmt.Sprintf("   📅 %s\n", bill.Date.Format("2006-01-02"))
		responseMsg += fmt.Sprintf("   📝 %d items\n\n", len(bill.Expenses))
	}

	return c.Send(responseMsg, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

func (h *BotHandler) handleSummaryBills(c tele.Context, userID string, intent *entities.Intent) error {
	bills, err := h.billWithExpensesService.ListBillsByUserID(userID)
	if err != nil {
		log.Printf("Failed to list bills: %v", err)
		return c.Send(h.messages.ErrorRetrieveBills)
	}

	if len(bills) == 0 {
		return c.Send(h.messages.NoBillsSummary)
	}

	// Determine period
	period := "all_time"
	if periodParam, ok := intent.Parameters["period"].(string); ok {
		period = periodParam
	}

	// Filter bills by period
	filteredBills := filterBillsByPeriod(bills, period)

	if len(filteredBills) == 0 {
		return c.Send(fmt.Sprintf(h.messages.NoBillsForPeriod, period))
	}

	// Calculate totals
	var totalPen, totalUsd float64
	categoryTotals := make(map[string]float64)

	for _, bill := range filteredBills {
		totalPen += bill.AmountPen
		totalUsd += bill.AmountUsd
		categoryTotals[bill.Category] += bill.AmountPen
	}

	// Build response
	periodName := getPeriodName(period)
	responseMsg := fmt.Sprintf(h.messages.SummaryHeader, periodName)
	responseMsg += fmt.Sprintf("💰 *Total Spent*\n")
	responseMsg += fmt.Sprintf("   PEN %.2f\n", totalPen)
	responseMsg += fmt.Sprintf("   USD %.2f\n\n", totalUsd)
	responseMsg += fmt.Sprintf("📋 *Number of Bills*: %d\n\n", len(filteredBills))

	if len(categoryTotals) > 0 {
		responseMsg += "*By Category (PEN)*:\n"
		for category, amount := range categoryTotals {
			responseMsg += fmt.Sprintf("   • %s: %.2f\n", category, amount)
		}
	}

	return c.Send(responseMsg, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

func getAmountInCurrency(bill *servicedtos.BillWithExpensesResponse, currency string) float64 {
	if currency == "PEN" {
		return bill.AmountPen
	}
	return bill.AmountUsd
}

func filterBillsByPeriod(bills []*servicedtos.BillWithExpensesResponse, period string) []*servicedtos.BillWithExpensesResponse {
	now := time.Now()
	var filtered []*servicedtos.BillWithExpensesResponse

	for _, bill := range bills {
		billDate := bill.Date

		switch period {
		case "last_month":
			lastMonth := now.AddDate(0, -1, 0)
			startOfLastMonth := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, lastMonth.Location())
			endOfLastMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)
			if billDate.After(startOfLastMonth) && billDate.Before(endOfLastMonth) {
				filtered = append(filtered, bill)
			}
		case "this_month":
			startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			if billDate.After(startOfMonth) {
				filtered = append(filtered, bill)
			}
		case "last_week":
			lastWeek := now.AddDate(0, 0, -7)
			if billDate.After(lastWeek) {
				filtered = append(filtered, bill)
			}
		case "all_time":
			fallthrough
		default:
			filtered = append(filtered, bill)
		}
	}

	return filtered
}

func getPeriodName(period string) string {
	switch period {
	case "last_month":
		return "Mes Pasado"
	case "this_month":
		return "Este Mes"
	case "last_week":
		return "Última Semana"
	case "all_time":
		return "Todo el Tiempo"
	default:
		return "Todo el Tiempo"
	}
}

func (h *BotHandler) handleCreateExpense(c tele.Context, userID string, intent *entities.Intent) error {
	// Extract parameters from intent
	amount, ok := intent.Parameters["amount"].(float64)
	if !ok || amount <= 0 {
		return c.Send(h.messages.ErrorMissingAmount)
	}

	// Get description from merchant or description parameter
	description := "Gasto"
	if merchant, ok := intent.Parameters["merchant"].(string); ok && merchant != "" {
		description = merchant
	} else if desc, ok := intent.Parameters["description"].(string); ok && desc != "" {
		description = desc
	}

	// Get category, default to "Other"
	category := "Other"
	if cat, ok := intent.Parameters["category"].(string); ok && cat != "" {
		category = cat
	}

	// Create bill with single expense
	now := time.Now()
	handlerDTO := handlerdtos.CreateBillWithExpensesRequest{
		UserID:       userID,
		Source:       "telegram",
		Description:  description,
		Category:     category,
		Currency:     "PEN", // Default to PEN for Peru
		ExchangeRate: 3.75,  // Default PEN to USD exchange rate
		Date:         now,
		Expenses: []handlerdtos.CreateExpenseForBill{
			{
				Description: description,
				Amount:      amount,
				Category:    category,
				Date:        now.Format("2006-01-02"),
			},
		},
	}

	serviceDTO := mappers.ToCreateBillWithExpensesServiceDTO(handlerDTO)
	bill, _, err := h.billWithExpensesService.CreateBillWithExpenses(serviceDTO)
	if err != nil {
		log.Printf("Failed to create expense: %v", err)
		return c.Send(h.messages.ErrorSaveExpense)
	}

	// Send success message
	responseMsg := fmt.Sprintf(h.messages.ExpenseSaved,
		amount,
		description,
		category,
		now.Format("2006-01-02"),
	)

	log.Printf("Expense created successfully: %s", bill.BillId)
	return c.Send(responseMsg, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

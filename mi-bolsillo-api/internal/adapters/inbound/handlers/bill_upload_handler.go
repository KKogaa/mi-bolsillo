package handlers

import (
	"io"
	"net/http"
	"time"

	"github.com/KKogaa/mi-bolsillo-api/internal/adapters/outbound/grok"
	handlerdtos "github.com/KKogaa/mi-bolsillo-api/internal/adapters/inbound/handlers/dtos"
	"github.com/KKogaa/mi-bolsillo-api/internal/adapters/inbound/handlers/mappers"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/services"
	"github.com/labstack/echo/v4"
)

type BillUploadHandler struct {
	grokClient              *grok.GrokClient
	billWithExpensesService *services.BillWithExpensesService
}

func NewBillUploadHandler(grokClient *grok.GrokClient, billWithExpensesService *services.BillWithExpensesService) *BillUploadHandler {
	return &BillUploadHandler{
		grokClient:              grokClient,
		billWithExpensesService: billWithExpensesService,
	}
}

// UploadBillPhoto godoc
// @Summary Upload a bill photo and parse it
// @Description Uploads a photo of a bill, parses it using Grok API, and creates a bill with expenses
// @Tags bills
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param image formData file true "Bill photo (JPEG, PNG, or other image format)"
// @Success 201 {object} map[string]interface{} "bill and expenses created successfully from image"
// @Failure 400 {object} map[string]string "Invalid request or image"
// @Failure 401 {object} map[string]string "User ID not found in context"
// @Failure 500 {object} map[string]string "Failed to process image or create bill"
// @Security BearerAuth
// @Router /bills/upload [post]
func (h *BillUploadHandler) UploadBillPhoto(c echo.Context) error {
	// Get user ID from context (set by Clerk auth middleware)
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User ID not found in context",
		})
	}

	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Image file is required",
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open uploaded file",
		})
	}
	defer src.Close()

	// Read file contents
	imageData, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to read uploaded file",
		})
	}

	// Parse the bill image using Grok
	parsedData, err := h.grokClient.ParseBillImage(imageData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse bill image: " + err.Error(),
		})
	}

	// Parse date string to time.Time
	billDate, err := time.Parse("2006-01-02", parsedData.Date)
	if err != nil {
		billDate = time.Now()
	}

	// Convert parsed data to CreateBillWithExpensesRequest
	handlerDTO := handlerdtos.CreateBillWithExpensesRequest{
		UserID:      userID,
		Description: parsedData.MerchantName,
		Category:    "General", // Default category for the bill
		Currency:    parsedData.Currency,
		Date:        billDate,
		Expenses:    make([]handlerdtos.CreateExpenseForBill, len(parsedData.Items)),
	}

	// Convert bill items to expenses
	for i, item := range parsedData.Items {
		handlerDTO.Expenses[i] = handlerdtos.CreateExpenseForBill{
			Description: item.Description,
			Amount:      item.Amount,
			Category:    item.Category,
			Date:        parsedData.Date,
		}
	}

	// Map handler DTO to service DTO using mapper
	serviceDTO := mappers.ToCreateBillWithExpensesServiceDTO(handlerDTO)

	// Create the bill with expenses
	bill, expenses, err := h.billWithExpensesService.CreateBillWithExpenses(serviceDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create bill with expenses: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"bill":        bill,
		"expenses":    expenses,
		"parsed_data": parsedData,
	})
}

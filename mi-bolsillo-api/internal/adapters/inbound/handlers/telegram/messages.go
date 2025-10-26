package telegram

import (
	"encoding/json"
	"fmt"
	"os"
)

// Messages holds all bot message templates
type Messages struct {
	Welcome            string `json:"welcome"`
	ProcessingImage    string `json:"processing_image"`
	BillSaved          string `json:"bill_saved"`
	ExpenseSaved       string `json:"expense_saved"`
	NoBills            string `json:"no_bills"`
	NoBillsSummary     string `json:"no_bills_summary"`
	NoBillsForPeriod   string `json:"no_bills_for_period"`
	BillsListHeader    string `json:"bills_list_header"`
	SummaryHeader      string `json:"summary_header"`
	UnknownIntent      string `json:"unknown_intent"`
	LinkAccountOTP     string `json:"link_account_otp"`
	LinkAccountError   string `json:"link_account_error"`
	ErrorUnderstand    string `json:"error_understand"`
	ErrorRetrieveImage string `json:"error_retrieve_image"`
	ErrorDownloadImage string `json:"error_download_image"`
	ErrorReadImage     string `json:"error_read_image"`
	ErrorParseBill     string `json:"error_parse_bill"`
	ErrorSaveBill      string `json:"error_save_bill"`
	ErrorSaveExpense   string `json:"error_save_expense"`
	ErrorRetrieveBills string `json:"error_retrieve_bills"`
	ErrorProcessingMsg string `json:"error_processing_message"`
	ErrorMissingAmount string `json:"error_missing_amount"`
}

// LoadMessages loads bot messages from a JSON file
func LoadMessages(path string) (*Messages, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read messages file: %w", err)
	}

	var messages Messages
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal messages: %w", err)
	}

	return &messages, nil
}

package entities

// Intent represents the detected user intent - a core domain entity
type Intent struct {
	Type       string                 `json:"type"`
	Confidence float64                `json:"confidence"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// Intent type constants - domain vocabulary
const (
	IntentListBills    = "list_bills"
	IntentSummaryBills = "summary_bills"
	IntentUploadBill   = "upload_bill"
	IntentUnknown      = "unknown"
)

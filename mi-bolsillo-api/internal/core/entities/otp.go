package entities

import "time"

// OTP represents a one-time password for account linking
type OTP struct {
	OTPCode    string    `json:"otpCode" db:"otp_code" example:"123456"`
	TelegramID int64     `json:"telegramId" db:"telegram_id" example:"123456789"`
	ExpiresAt  time.Time `json:"expiresAt" db:"expires_at" example:"2025-10-10T10:05:00Z"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at" example:"2025-10-10T10:00:00Z"`
}

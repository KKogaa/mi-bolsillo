package entities

import "time"

// User represents a user entity that can have both Telegram and Clerk identities
type User struct {
	UserID     string    `json:"userId" db:"user_id" example:"user_123456789"`
	ClerkID    *string   `json:"clerkId,omitempty" db:"clerk_id" example:"user_2abc123def456"`
	TelegramID *int64    `json:"telegramId,omitempty" db:"telegram_id" example:"123456789"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at" example:"2025-10-10T10:00:00Z"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at" example:"2025-10-10T10:00:00Z"`
}

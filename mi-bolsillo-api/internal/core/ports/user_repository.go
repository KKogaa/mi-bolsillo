package ports

import "github.com/KKogaa/mi-bolsillo-api/internal/core/entities"

type UserRepository interface {
	Create(user *entities.User) error
	FindByID(userID string) (*entities.User, error)
	FindByClerkID(clerkID string) (*entities.User, error)
	FindByTelegramID(telegramID int64) (*entities.User, error)
	Update(user *entities.User) error
	LinkClerkAccount(userID string, clerkID string) error
}

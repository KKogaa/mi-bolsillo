package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/KKogaa/mi-bolsillo-api/internal/core/entities"
	"github.com/KKogaa/mi-bolsillo-api/internal/core/ports"
	"github.com/google/uuid"
)

type AccountLinkService struct {
	userRepo    ports.UserRepository
	otpRepo     ports.OTPRepository
	billRepo    ports.BillRepository
	expenseRepo ports.ExpenseRepository
	otpExpirationMinutes int
}

func NewAccountLinkService(
	userRepo ports.UserRepository,
	otpRepo ports.OTPRepository,
	billRepo ports.BillRepository,
	expenseRepo ports.ExpenseRepository,
	otpExpirationMinutes int,
) *AccountLinkService {
	return &AccountLinkService{
		userRepo:    userRepo,
		otpRepo:     otpRepo,
		billRepo:    billRepo,
		expenseRepo: expenseRepo,
		otpExpirationMinutes: otpExpirationMinutes,
	}
}

// GenerateOTP creates a new OTP for the given Telegram user
func (s *AccountLinkService) GenerateOTP(telegramID int64) (string, error) {
	// Clean up expired OTPs first
	_ = s.otpRepo.DeleteExpired()

	// Generate a random 6-digit OTP
	otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))

	otp := &entities.OTP{
		OTPCode:    otpCode,
		TelegramID: telegramID,
		ExpiresAt:  time.Now().Add(time.Duration(s.otpExpirationMinutes) * time.Minute),
		CreatedAt:  time.Now(),
	}

	if err := s.otpRepo.Create(otp); err != nil {
		return "", fmt.Errorf("failed to create OTP: %w", err)
	}

	return otpCode, nil
}

// VerifyAndLinkAccounts validates the OTP and links the Telegram account with the Clerk account
func (s *AccountLinkService) VerifyAndLinkAccounts(otpCode string, clerkID string) error {
	// Find the OTP
	otp, err := s.otpRepo.FindByCode(otpCode)
	if err != nil {
		return fmt.Errorf("failed to find OTP: %w", err)
	}

	if otp == nil {
		return ErrInvalidOTP
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		_ = s.otpRepo.Delete(otpCode)
		return ErrOTPExpired
	}

	// Check if Clerk user already exists
	existingClerkUser, err := s.userRepo.FindByClerkID(clerkID)
	if err != nil {
		return fmt.Errorf("failed to find clerk user: %w", err)
	}

	// Check if Telegram user already exists
	existingTelegramUser, err := s.userRepo.FindByTelegramID(otp.TelegramID)
	if err != nil {
		return fmt.Errorf("failed to find telegram user: %w", err)
	}

	now := time.Now()

	// Case 1: Both users exist - merge them
	if existingClerkUser != nil && existingTelegramUser != nil {
		// If they're already the same user, just update
		if existingClerkUser.UserID == existingTelegramUser.UserID {
			_ = s.otpRepo.Delete(otpCode)
			return nil
		}

		// Prioritize the Clerk user (existing web user) and add Telegram ID to it
		existingClerkUser.TelegramID = &otp.TelegramID
		existingClerkUser.UpdatedAt = now
		if err := s.userRepo.Update(existingClerkUser); err != nil {
			return fmt.Errorf("failed to update clerk user with telegram id: %w", err)
		}

		// Migrate bills from Telegram user to Clerk user
		if err := s.migrateBills(existingTelegramUser.UserID, existingClerkUser.UserID); err != nil {
			return fmt.Errorf("failed to migrate bills from telegram user: %w", err)
		}

		_ = s.otpRepo.Delete(otpCode)
		return nil
	}

	// Case 2: Only Clerk user exists - add Telegram ID
	if existingClerkUser != nil {
		existingClerkUser.TelegramID = &otp.TelegramID
		existingClerkUser.UpdatedAt = now
		if err := s.userRepo.Update(existingClerkUser); err != nil {
			return fmt.Errorf("failed to update clerk user: %w", err)
		}

		_ = s.otpRepo.Delete(otpCode)
		return nil
	}

	// Case 3: Only Telegram user exists - add Clerk ID
	if existingTelegramUser != nil {
		existingTelegramUser.ClerkID = &clerkID
		existingTelegramUser.UpdatedAt = now
		if err := s.userRepo.Update(existingTelegramUser); err != nil {
			return fmt.Errorf("failed to update telegram user: %w", err)
		}

		_ = s.otpRepo.Delete(otpCode)
		return nil
	}

	// Case 4: Neither user exists - create a new linked user
	newUser := &entities.User{
		UserID:     uuid.New().String(),
		ClerkID:    &clerkID,
		TelegramID: &otp.TelegramID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	_ = s.otpRepo.Delete(otpCode)
	return nil
}

// GetUserByTelegramID retrieves a user by their Telegram ID
func (s *AccountLinkService) GetUserByTelegramID(telegramID int64) (*entities.User, error) {
	return s.userRepo.FindByTelegramID(telegramID)
}

// GetUserByClerkID retrieves a user by their Clerk ID
func (s *AccountLinkService) GetUserByClerkID(clerkID string) (*entities.User, error) {
	return s.userRepo.FindByClerkID(clerkID)
}

// GetOrCreateUserByTelegramID gets an existing user or creates a new one for a Telegram user
func (s *AccountLinkService) GetOrCreateUserByTelegramID(telegramID int64) (*entities.User, error) {
	user, err := s.userRepo.FindByTelegramID(telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to find telegram user: %w", err)
	}

	if user != nil {
		return user, nil
	}

	// Create new user
	now := time.Now()
	newUser := &entities.User{
		UserID:     uuid.New().String(),
		TelegramID: &telegramID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

// GetOrCreateUserByClerkID gets an existing user or creates a new one for a Clerk user
func (s *AccountLinkService) GetOrCreateUserByClerkID(clerkID string) (*entities.User, error) {
	user, err := s.userRepo.FindByClerkID(clerkID)
	if err != nil {
		return nil, fmt.Errorf("failed to find clerk user: %w", err)
	}

	if user != nil {
		return user, nil
	}

	// Create new user
	now := time.Now()
	newUser := &entities.User{
		UserID:    uuid.New().String(),
		ClerkID:   &clerkID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

// migrateBills updates all bills and expenses from oldUserID to newUserID
func (s *AccountLinkService) migrateBills(oldUserID string, newUserID string) error {
	// Update bills
	if err := s.billRepo.UpdateUserID(oldUserID, newUserID); err != nil {
		return fmt.Errorf("failed to migrate bills: %w", err)
	}

	// Update expenses
	if err := s.expenseRepo.UpdateUserID(oldUserID, newUserID); err != nil {
		return fmt.Errorf("failed to migrate expenses: %w", err)
	}

	return nil
}

var (
	ErrInvalidOTP = errors.New("invalid OTP code")
	ErrOTPExpired = errors.New("OTP has expired")
)

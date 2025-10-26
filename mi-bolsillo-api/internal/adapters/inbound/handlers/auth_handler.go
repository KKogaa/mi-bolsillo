package handlers

import (
	"net/http"

	"github.com/KKogaa/mi-bolsillo-api/internal/core/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	accountLinkService *services.AccountLinkService
}

func NewAuthHandler(accountLinkService *services.AccountLinkService) *AuthHandler {
	return &AuthHandler{
		accountLinkService: accountLinkService,
	}
}

type VerifyOTPRequest struct {
	OTPCode string `json:"otpCode" example:"123456"`
}

type VerifyOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LinkStatusResponse struct {
	IsLinked   bool   `json:"isLinked"`
	TelegramID *int64 `json:"telegramId,omitempty"`
}

// VerifyOTP godoc
// @Summary Verify OTP and link Telegram account to Clerk account
// @Description Validates the OTP code generated from Telegram and links it to the authenticated Clerk user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "OTP Code"
// @Success 200 {object} VerifyOTPResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c echo.Context) error {
	// Get user ID from JWT (set by Clerk middleware)
	clerkID, ok := c.Get("userID").(string)
	if !ok || clerkID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized - invalid user ID",
		})
	}

	var req VerifyOTPRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if req.OTPCode == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "OTP code is required",
		})
	}

	// Verify OTP and link accounts
	if err := h.accountLinkService.VerifyAndLinkAccounts(req.OTPCode, clerkID); err != nil {
		if err == services.ErrInvalidOTP {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid OTP code",
			})
		}
		if err == services.ErrOTPExpired {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "OTP code has expired",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to verify OTP and link accounts",
		})
	}

	return c.JSON(http.StatusOK, VerifyOTPResponse{
		Success: true,
		Message: "Accounts successfully linked",
	})
}

// GetLinkStatus godoc
// @Summary Check if account is linked to Telegram
// @Description Returns whether the authenticated Clerk account is linked to a Telegram account
// @Tags auth
// @Produce json
// @Success 200 {object} LinkStatusResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /auth/link-status [get]
func (h *AuthHandler) GetLinkStatus(c echo.Context) error {
	// Get user ID from JWT (set by Clerk middleware)
	clerkID, ok := c.Get("userID").(string)
	if !ok || clerkID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized - invalid user ID",
		})
	}

	// Get user by Clerk ID
	user, err := h.accountLinkService.GetUserByClerkID(clerkID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user information",
		})
	}

	// If user doesn't exist or has no Telegram ID, they're not linked
	if user == nil || user.TelegramID == nil {
		return c.JSON(http.StatusOK, LinkStatusResponse{
			IsLinked:   false,
			TelegramID: nil,
		})
	}

	// User is linked
	return c.JSON(http.StatusOK, LinkStatusResponse{
		IsLinked:   true,
		TelegramID: user.TelegramID,
	})
}

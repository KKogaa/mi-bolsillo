package ports

import "github.com/KKogaa/mi-bolsillo-api/internal/core/entities"

type OTPRepository interface {
	Create(otp *entities.OTP) error
	FindByCode(otpCode string) (*entities.OTP, error)
	Delete(otpCode string) error
	DeleteExpired() error
}

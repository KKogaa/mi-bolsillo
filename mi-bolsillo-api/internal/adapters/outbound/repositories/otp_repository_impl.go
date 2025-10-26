package repositories

import (
	"database/sql"
	"errors"

	"github.com/KKogaa/mi-bolsillo-api/internal/core/entities"
	"github.com/jmoiron/sqlx"
)

type OTPRepositoryImpl struct {
	db *sqlx.DB
}

func NewOTPRepository(db *sqlx.DB) *OTPRepositoryImpl {
	return &OTPRepositoryImpl{db: db}
}

func (r *OTPRepositoryImpl) Create(otp *entities.OTP) error {
	query := `
		INSERT INTO account_link_otps (otp_code, telegram_id, expires_at, created_at)
		VALUES (:otp_code, :telegram_id, :expires_at, :created_at)
	`
	_, err := r.db.NamedExec(query, otp)
	return err
}

func (r *OTPRepositoryImpl) FindByCode(otpCode string) (*entities.OTP, error) {
	var otp entities.OTP
	query := `SELECT * FROM account_link_otps WHERE otp_code = ?`
	err := r.db.Get(&otp, query, otpCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &otp, nil
}

func (r *OTPRepositoryImpl) Delete(otpCode string) error {
	query := `DELETE FROM account_link_otps WHERE otp_code = ?`
	_, err := r.db.Exec(query, otpCode)
	return err
}

func (r *OTPRepositoryImpl) DeleteExpired() error {
	query := `DELETE FROM account_link_otps WHERE expires_at < CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query)
	return err
}

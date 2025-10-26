package repositories

import (
	"database/sql"
	"errors"

	"github.com/KKogaa/mi-bolsillo-api/internal/core/entities"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(user *entities.User) error {
	query := `
		INSERT INTO users (user_id, clerk_id, telegram_id, created_at, updated_at)
		VALUES (:user_id, :clerk_id, :telegram_id, :created_at, :updated_at)
	`
	_, err := r.db.NamedExec(query, user)
	return err
}

func (r *UserRepositoryImpl) FindByID(userID string) (*entities.User, error) {
	var user entities.User
	query := `SELECT * FROM users WHERE user_id = ?`
	err := r.db.Get(&user, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByClerkID(clerkID string) (*entities.User, error) {
	var user entities.User
	query := `SELECT * FROM users WHERE clerk_id = ?`
	err := r.db.Get(&user, query, clerkID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByTelegramID(telegramID int64) (*entities.User, error) {
	var user entities.User
	query := `SELECT * FROM users WHERE telegram_id = ?`
	err := r.db.Get(&user, query, telegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) Update(user *entities.User) error {
	query := `
		UPDATE users
		SET clerk_id = :clerk_id, telegram_id = :telegram_id, updated_at = :updated_at
		WHERE user_id = :user_id
	`
	_, err := r.db.NamedExec(query, user)
	return err
}

func (r *UserRepositoryImpl) LinkClerkAccount(userID string, clerkID string) error {
	query := `
		UPDATE users
		SET clerk_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = ?
	`
	_, err := r.db.Exec(query, clerkID, userID)
	return err
}

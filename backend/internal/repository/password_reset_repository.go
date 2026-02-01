package repository

import (
	"auth-go/internal/domain"

	"gorm.io/gorm"
)

type passwordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) domain.PasswordResetRepository {
	return &passwordResetRepository{db}
}

func (r *passwordResetRepository) Save(reset *domain.PasswordResetToken) (*domain.PasswordResetToken, error) {
	err := r.db.Create(reset).Error
	if err != nil {
		return nil, err
	}
	return reset, nil
}

func (r *passwordResetRepository) FindByToken(token string) (*domain.PasswordResetToken, error) {
	var reset domain.PasswordResetToken
	err := r.db.Where("token = ?", token).First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

func (r *passwordResetRepository) DeleteByEmail(email string) error {
	return r.db.Where("email = ?", email).Delete(&domain.PasswordResetToken{}).Error
}

package service

import (
	"auth-go/internal/config"
	"auth-go/internal/domain"
	"auth-go/pkg/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(input *domain.RegisterInput) (*domain.User, error)
	Login(input *domain.LoginInput) (string, *domain.User, error)
	ForgotPassword(input *domain.ForgotPasswordInput) error
	ResetPassword(input *domain.ResetPasswordInput) error
}

type authService struct {
	userRepo     domain.UserRepository
	resetRepo    domain.PasswordResetRepository
	emailService EmailService
	config       *config.Config
}

func NewAuthService(userRepo domain.UserRepository, resetRepo domain.PasswordResetRepository, emailService EmailService, config *config.Config) AuthService {
	return &authService{userRepo, resetRepo, emailService, config}
}

func (s *authService) Register(input *domain.RegisterInput) (*domain.User, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.FindByEmail(input.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	savedUser, err := s.userRepo.Save(newUser)
	if err != nil {
		return nil, err
	}

	// Send welcome email (async to avoid blocking)
	go s.emailService.SendWelcomeEmail(savedUser.Email, savedUser.Name)

	return savedUser, nil
}

func (s *authService) Login(input *domain.LoginInput) (string, *domain.User, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid email or password")
		}
		return "", nil, err
	}

	// Check password
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return "", nil, errors.New("invalid email or password")
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, user.Email, s.config.JWTSecret, s.config.JWTExpiredIn)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *authService) ForgotPassword(input *domain.ForgotPasswordInput) error {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		// Return nil to avoid email enumeration
		return nil
	}

	// Generate token (simple random string for now, better to use crypto rand)
	resetToken, _ := utils.HashPassword(time.Now().String() + user.Email) // Simple hack for unique string

	// Save token
	resetData := &domain.PasswordResetToken{
		Email:     user.Email,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// Remove old tokens
	s.resetRepo.DeleteByEmail(user.Email)

	_, err = s.resetRepo.Save(resetData)
	if err != nil {
		return err
	}

	// Send email
	resetLink := fmt.Sprintf("http://localhost:5173/reset-password?token=%s&email=%s", resetToken, user.Email)
	go s.emailService.SendResetPasswordEmail(user.Email, resetLink)

	return nil
}

func (s *authService) ResetPassword(input *domain.ResetPasswordInput) error {
	// Validate token
	resetData, err := s.resetRepo.FindByToken(input.Token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if time.Now().After(resetData.ExpiresAt) {
		return errors.New("token expired")
	}

	// Update user password
	user, err := s.userRepo.FindByEmail(resetData.Email)
	if err != nil {
		return errors.New("user not found")
	}

	hashedPassword, _ := utils.HashPassword(input.Password)
	user.Password = hashedPassword

	_, err = s.userRepo.Update(user)
	if err != nil {
		return err
	}

	// Delete used token
	s.resetRepo.DeleteByEmail(user.Email)

	return nil
}

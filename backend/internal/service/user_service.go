package service

import (
	"auth-go/internal/domain"
	"errors"
)

type UserService interface {
	GetProfile(userID uint64) (*domain.User, error)
	GetAllUsers(page int, limit int, search string) ([]*domain.User, int64, error)
}

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) GetProfile(userID uint64) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	// Sanitize output just in case (e.g. remove password)
	user.Password = ""
	return user, nil
}

func (s *userService) GetAllUsers(page int, limit int, search string) ([]*domain.User, int64, error) {
	users, total, err := s.userRepo.FindAll(page, limit, search)
	if err != nil {
		return nil, 0, err
	}

	// Sanitize passwords
	for _, user := range users {
		user.Password = ""
	}

	return users, total, nil
}

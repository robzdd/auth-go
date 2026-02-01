package repository

import (
	"auth-go/internal/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Save(user *domain.User) (*domain.User, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint64) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(page int, limit int, search string) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	// Base query
	query := r.db.Model(&domain.User{})

	// Search filter
	if search != "" {
		// Optimization: Use separate queries per index or use just one if performance is critical for "OR"
		// LIKE 'val%' uses index. LIKE '%val%' does NOT.
		// For 3 million rows, we MUST use Prefix match ('val%') or proper FullText search.
		// using "OR" with two columns can sometimes skip index utilization depending on MySQL version.
		// For now, let's optimize to prefix match on Name OR Email (email already indexed).
		searchPattern := search + "%"
		query = query.Where("name LIKE ? OR email LIKE ?", searchPattern, searchPattern)
	}

	// Count total records (before pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(user *domain.User) (*domain.User, error) {
	err := r.db.Save(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

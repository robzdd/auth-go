package main

import (
	"fmt"
	"log"
	"time"

	"auth-go/internal/config"
	"auth-go/internal/database"
	"auth-go/internal/domain"
	"auth-go/pkg/utils"

	"gorm.io/gorm"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Connect Database
	db := database.ConnectDB(cfg)

	// 3. Seed Data
	log.Println("Seeding database...")
	seedUsers(db)
	log.Println("Database seeded successfully!")
}

func seedUsers(db *gorm.DB) {
	// 1. Hash password once (very important for performance)
	password, err := utils.HashPassword("password")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	totalUsers := 3000000
	batchSize := 2000 // Insert 2000 users per query

	// Create "Admin User" specifically first if not exists
	var admin domain.User
	if err := db.Where("email = ?", "admin@example.com").First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			db.Create(&domain.User{
				Name:      "Admin User",
				Email:     "admin@example.com",
				Password:  password,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
			log.Println("Created Admin User")
		}
	}

	start := time.Now()
	users := make([]domain.User, 0, batchSize)

	// Start from 1 up to totalUsers
	count := 0

	// Transaction optimization can also help, but GORM batch insert wraps it usually.
	// We'll manually chunk the loop to manage memory.

	log.Printf("Starting bulk insert of %d users...", totalUsers)

	for i := 1; i <= totalUsers; i++ {
		users = append(users, domain.User{
			Name:      fmt.Sprintf("User %d", i),
			Email:     fmt.Sprintf("user%d@example.com", i),
			Password:  password, // Reuse hashed password
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		// If batch is full, insert and reset
		if len(users) >= batchSize {
			if err := db.CreateInBatches(users, batchSize).Error; err != nil {
				log.Printf("Error inserting batch at index %d: %v", i, err)
			}
			users = users[:0] // Reset slice but keep capacity
			count += batchSize

			if count%100000 == 0 {
				log.Printf("Inserted %d users... (Elapsed: %v)", count, time.Since(start))
			}
		}
	}

	// Insert potential remaining users
	if len(users) > 0 {
		db.CreateInBatches(users, len(users))
		count += len(users)
	}

	log.Printf("Finished! Total inserted: %d users in %v", count, time.Since(start))
}

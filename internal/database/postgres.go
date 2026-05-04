package database

import (
	"fmt"

	"github.com/user/go-commerce-api/internal/config"
	"github.com/user/go-commerce-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migration
	err = db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Product{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

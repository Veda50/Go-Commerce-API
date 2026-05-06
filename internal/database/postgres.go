package database

import (
	"fmt"

	"github.com/user/go-commerce-api/internal/config"
	"github.com/user/go-commerce-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	if cfg.DBURL != "" {
		dialector = postgres.Open(cfg.DBURL)
	} else {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
		dialector = postgres.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migration
	err = db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.Cart{},
		&model.CartItem{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

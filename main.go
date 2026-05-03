package main

import (
	"fmt"
	"log"

	"github.com/user/go-commerce-api/internal/config"
	"github.com/user/go-commerce-api/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	fmt.Printf("Connected to database %s on %s:%s\n", cfg.DBName, cfg.DBHost, cfg.DBPort)

	// To keep the process running for now
	_ = db
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/user/go-commerce-api/internal/config"
	"github.com/user/go-commerce-api/internal/database"
	"github.com/user/go-commerce-api/internal/handler"
	"github.com/user/go-commerce-api/internal/middleware"
	"github.com/user/go-commerce-api/internal/repository"
	"github.com/user/go-commerce-api/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	fmt.Println("Database connection and migration successful")

	// Repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepo, productRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	adminHandler := handler.NewAdminHandler(productService)
	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)

	// Router
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.Signup)
		r.Post("/login", authHandler.Login)
	})

	r.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.ListProducts)
		r.Get("/{id}", productHandler.GetProduct)
	})

	r.Route("/cart", func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTSecret))
		r.Get("/", cartHandler.GetCart)
		r.Post("/items", cartHandler.AddToCart)
		r.Delete("/items/{id}", cartHandler.RemoveItem)
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTSecret))
		r.Use(middleware.AdminOnly)
		r.Post("/products", adminHandler.CreateProduct)
		r.Put("/products/{id}", adminHandler.UpdateProduct)
		r.Delete("/products/{id}", adminHandler.DeleteProduct)
	})

	fmt.Printf("Server starting on port %s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}

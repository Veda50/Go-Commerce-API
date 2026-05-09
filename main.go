package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/user/go-commerce-api/src/config"
	"github.com/user/go-commerce-api/src/database"
	"github.com/user/go-commerce-api/src/handler"
	"github.com/user/go-commerce-api/src/middleware"
	"github.com/user/go-commerce-api/src/repository"
	"github.com/user/go-commerce-api/src/service"
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
	orderRepo := repository.NewOrderRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	paymentService := service.NewPaymentService(cfg.XenditKey)
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo, userRepo, paymentService, db)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	adminHandler := handler.NewAdminHandler(productService, orderRepo)
	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService, orderRepo)
	webhookHandler := handler.NewWebhookHandler(orderService, cfg)

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
		r.Patch("/items/{id}", cartHandler.UpdateQuantity)
		r.Delete("/items/{id}", cartHandler.RemoveItem)
	})

	r.Route("/orders", func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTSecret))
		r.Get("/", orderHandler.GetHistory)
	})

	r.Post("/checkout", func(w http.ResponseWriter, r *http.Request) {
		middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(orderHandler.Checkout)).ServeHTTP(w, r)
	})

	r.Post("/webhook/xendit", webhookHandler.XenditCallback)

	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTSecret))
		r.Use(middleware.AdminOnly)
		r.Post("/products", adminHandler.CreateProduct)
		r.Put("/products/{id}", adminHandler.UpdateProduct)
		r.Delete("/products/{id}", adminHandler.DeleteProduct)
		r.Patch("/products/{id}/stock", adminHandler.UpdateStock)
		r.Get("/orders", adminHandler.ListOrders)
	})

	fmt.Printf("Server starting on port %s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}

# Go E-Commerce API

Modular RESTful API for E-Commerce built with Go, GORM, and PostgreSQL.

## Portfolio
This project is built to maintain and demonstrate proficiency in Go backend development, focusing on modular architecture, database transactions, and payment gateway integration.

## Features

### Completed
- JWT Authentication and RBAC
- Product and Category CRUD
- Cart management with stock validation
- Transactional checkout with Xendit Invoice
- Webhook callback for payment status
- Admin stock management and order tracking
- Centralized JSON response handling

### Soon
- S3/R2 Image uploads
- Product reviews and ratings
- Shipping API integration
- Redis caching
- Unit testing

## Tech Stack
- Go, Chi, GORM, PostgreSQL, Xendit

## Setup

1. Copy .env.example to .env and fill variables
2. go mod download
3. Expose port 8080 for webhooks (use ngrok, Xendit CLI, or my personal `tunnel.vedabe.com` subdomain)
4. go run main.go

## API Endpoints

### Auth
- POST /auth/signup
- POST /auth/login

### Products
- GET /products
- GET /products/{id}

### Cart
- GET /cart
- POST /cart/items
- PATCH /cart/items/{id}
- DELETE /cart/items/{id}

### Orders
- POST /checkout
- GET /orders

### Admin
- POST /admin/products
- PUT /admin/products/{id}
- DELETE /admin/products/{id}
- PATCH /admin/products/{id}/stock
- GET /admin/orders

### Webhook
- POST /webhook/xendit

# Go E-Commerce API

A modular RESTful API for E-Commerce built with Go, GORM, and PostgreSQL.

## Features
- **Auth**: Signup & Login with JWT.
- **Products**: CRUD for products and categories (Admin only for CUD).
- **Cart**: Manage cart items, stock validation, and subtotal calculation.
- **Orders**: Transactional checkout process with Xendit Invoice integration.
- **Admin**: Update stock, view all orders with status filtering.
- **Payments**: Xendit Invoice with Webhook callback support.

## Tech Stack
- **Language**: Go 1.21+
- **Framework**: `go-chi/chi` (Router)
- **ORM**: `GORM` with PostgreSQL driver
- **Auth**: `golang-jwt/jwt`, `bcrypt`
- **Payment**: `xendit-go/v6`

## Setup Guide

1. **Clone the repository**
   ```bash
   git clone https://github.com/user/go-commerce-api.git
   cd go-commerce-api
   ```

2. **Configure Environment Variables**
   Copy `.env.example` to `.env` and fill in the values.
   ```bash
   cp .env.example .env
   ```

3. **Install Dependencies**
   ```bash
   go mod download
   ```

4. **Run the Application**
   ```bash
   go run main.go
   ```

## API Endpoints

### Auth
- `POST /auth/signup` - Register new user
- `POST /auth/login` - Login and get JWT

### Products
- `GET /products` - List products (query param `q` for search)
- `GET /products/{id}` - Get product detail

### Cart (Protected)
- `GET /cart` - View current active cart
- `POST /cart/items` - Add item to cart
- `PATCH /cart/items/{id}` - Update item quantity
- `DELETE /cart/items/{id}` - Remove item from cart

### Orders (Protected)
- `POST /checkout` - Create order from cart and get Xendit Invoice URL
- `GET /orders` - View order history

### Admin (Admin Only)
- `POST /admin/products` - Create product
- `PUT /admin/products/{id}` - Update product
- `DELETE /admin/products/{id}` - Delete product
- `PATCH /admin/products/{id}/stock` - Update stock manually
- `GET /admin/orders` - View all orders (filter `status=PAID` etc)

### Webhook
- `POST /webhook/xendit` - Xendit payment callback

---
sidebar_position: 3
---

# Microservices Example

Learn how to build a microservices architecture with go-infra, including service discovery, inter-service communication, and distributed tracing.

## Overview

This example demonstrates a complete microservices setup with:

- **User Service** - User management and authentication
- **Product Service** - Product catalog management
- **Order Service** - Order processing and fulfillment
- **API Gateway** - Unified entry point
- **Service Discovery** - Dynamic service registration
- **Message Queue** - Async communication with RabbitMQ
- **Distributed Tracing** - OpenTelemetry integration

## Architecture

```
                    ┌─────────────────┐
                    │   API Gateway   │
                    │   (Port 8000)   │
                    └────────┬────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
    ┌────▼─────┐      ┌─────▼──────┐     ┌─────▼──────┐
    │  User    │      │  Product   │     │   Order    │
    │ Service  │      │  Service   │     │  Service   │
    │(Port 8001)│     │(Port 8002)│     │(Port 8003)│
    └────┬─────┘      └─────┬──────┘     └─────┬──────┘
         │                   │                   │
         └───────────────────┼───────────────────┘
                             │
                    ┌────────▼────────┐
                    │   PostgreSQL    │
                    │    (Shared)     │
                    └─────────────────┘
```

## Project Structure

```
examples/microservices/
├── api-gateway/
│   ├── main.go
│   ├── routes.go
│   └── middleware/
├── user-service/
│   ├── main.go
│   ├── handler/
│   ├── service/
│   └── repository/
├── product-service/
│   ├── main.go
│   ├── handler/
│   ├── service/
│   └── repository/
├── order-service/
│   ├── main.go
│   ├── handler/
│   ├── service/
│   └── repository/
├── shared/
│   ├── events/
│   ├── models/
│   └── client/
├── docker-compose.yml
└── README.md
```

## User Service

### Main Server

```go
// user-service/main.go
package main

import (
    "log"
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber"
    "github.com/phatnt199/go-infra/pkg/adapter/http/crud"
    "github.com/phatnt199/go-infra/pkg/logger"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type User struct {
    ID    string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email" gorm:"uniqueIndex"`
    Role  string `json:"role" default:"user"`
}

func main() {
    logger.Init()

    // Database connection
    dsn := "host=localhost port=5432 user=postgres password=postgres dbname=users"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    db.AutoMigrate(&User{})

    // HTTP server
    app, router := fiber.New()

    // CRUD endpoints
    crud.RegisterCRUD[*User](router, db, &crud.CRUDOptions[*User, string]{
        BasePath: "/api/users",
    })

    // Health check
    router.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "healthy"})
    })

    log.Println("User service starting on :8001")
    log.Fatal(app.Listen(":8001"))
}
```

## Product Service

```go
// product-service/main.go
package main

import (
    "log"
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber"
    "github.com/phatnt199/go-infra/pkg/adapter/http/crud"
    "github.com/phatnt199/go-infra/pkg/logger"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Product struct {
    ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Name        string  `json:"name" validate:"required"`
    Description string  `json:"description"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Stock       int     `json:"stock" validate:"gte=0"`
}

func main() {
    logger.Init()

    // Database connection
    dsn := "host=localhost port=5432 user=postgres password=postgres dbname=products"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    db.AutoMigrate(&Product{})

    // HTTP server
    app, router := fiber.New()

    // CRUD endpoints
    crud.RegisterCRUD[*Product](router, db, &crud.CRUDOptions[*Product, string]{
        BasePath: "/api/products",
    })

    // Health check
    router.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "healthy"})
    })

    log.Println("Product service starting on :8002")
    log.Fatal(app.Listen(":8002"))
}
```

## Order Service

### Order Model

```go
// order-service/models/order.go
package models

import "time"

type Order struct {
    ID          string      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID      string      `json:"user_id" validate:"required"`
    Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
    TotalAmount float64     `json:"total_amount"`
    Status      string      `json:"status" default:"pending"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
    ID        string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    OrderID   string  `json:"order_id"`
    ProductID string  `json:"product_id" validate:"required"`
    Quantity  int     `json:"quantity" validate:"required,gt=0"`
    Price     float64 `json:"price"`
}
```

### Order Service

```go
// order-service/service/order_service.go
package service

import (
    "emperror.dev/errors"
    "github.com/phatnt199/go-infra/pkg/logger"
    "order-service/models"
    "order-service/client"
)

type OrderService struct {
    productClient *client.ProductClient
    userClient    *client.UserClient
    logger        logger.Logger
}

func NewOrderService(
    productClient *client.ProductClient,
    userClient *client.UserClient,
    logger logger.Logger,
) *OrderService {
    return &OrderService{
        productClient: productClient,
        userClient:    userClient,
        logger:        logger,
    }
}

func (s *OrderService) CreateOrder(order *models.Order) error {
    // Verify user exists
    user, err := s.userClient.GetUser(order.UserID)
    if err != nil {
        return errors.Wrap(err, "user not found")
    }

    s.logger.Info("Creating order for user",
        logger.Field("userId", user.ID),
        logger.Field("userName", user.Name))

    // Verify products and calculate total
    var totalAmount float64
    for i, item := range order.Items {
        product, err := s.productClient.GetProduct(item.ProductID)
        if err != nil {
            return errors.Wrap(err, "product not found")
        }

        // Check stock
        if product.Stock < item.Quantity {
            return errors.New("insufficient stock")
        }

        // Set price from product
        order.Items[i].Price = product.Price
        totalAmount += product.Price * float64(item.Quantity)
    }

    order.TotalAmount = totalAmount
    order.Status = "pending"

    s.logger.Info("Order created successfully",
        logger.Field("orderId", order.ID),
        logger.Field("totalAmount", totalAmount))

    return nil
}
```

### HTTP Client

```go
// order-service/client/product_client.go
package client

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Product struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
}

type ProductClient struct {
    baseURL string
    client  *http.Client
}

func NewProductClient(baseURL string) *ProductClient {
    return &ProductClient{
        baseURL: baseURL,
        client:  &http.Client{},
    }
}

func (c *ProductClient) GetProduct(id string) (*Product, error) {
    url := fmt.Sprintf("%s/api/products/%s", c.baseURL, id)

    resp, err := c.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("product not found")
    }

    var product Product
    if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
        return nil, err
    }

    return &product, nil
}
```

## API Gateway

```go
// api-gateway/main.go
package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
    app := fiber.New()

    // Middleware
    app.Use(logger.New())
    app.Use(cors.New())

    // Service URLs
    userServiceURL, _ := url.Parse("http://localhost:8001")
    productServiceURL, _ := url.Parse("http://localhost:8002")
    orderServiceURL, _ := url.Parse("http://localhost:8003")

    // Proxy to services
    app.All("/api/users/*", proxyHandler(userServiceURL))
    app.All("/api/products/*", proxyHandler(productServiceURL))
    app.All("/api/orders/*", proxyHandler(orderServiceURL))

    // Health check
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "healthy",
            "services": map[string]string{
                "user":    "http://localhost:8001",
                "product": "http://localhost:8002",
                "order":   "http://localhost:8003",
            },
        })
    })

    log.Println("API Gateway starting on :8000")
    log.Fatal(app.Listen(":8000"))
}

func proxyHandler(target *url.URL) fiber.Handler {
    proxy := httputil.NewSingleHostReverseProxy(target)

    return func(c *fiber.Ctx) error {
        proxy.ServeHTTP(c.Response().BodyWriter(), c.Request())
        return nil
    }
}
```

## Docker Compose

```yaml
# docker-compose.yml
version: "3.8"

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  user-service:
    build: ./user-service
    ports:
      - "8001:8001"
    environment:
      DATABASE_URL: postgres://postgres:postgres@postgres:5432/users
    depends_on:
      - postgres

  product-service:
    build: ./product-service
    ports:
      - "8002:8002"
    environment:
      DATABASE_URL: postgres://postgres:postgres@postgres:5432/products
    depends_on:
      - postgres

  order-service:
    build: ./order-service
    ports:
      - "8003:8003"
    environment:
      DATABASE_URL: postgres://postgres:postgres@postgres:5432/orders
      USER_SERVICE_URL: http://user-service:8001
      PRODUCT_SERVICE_URL: http://product-service:8002
    depends_on:
      - postgres
      - user-service
      - product-service

  api-gateway:
    build: ./api-gateway
    ports:
      - "8000:8000"
    environment:
      USER_SERVICE_URL: http://user-service:8001
      PRODUCT_SERVICE_URL: http://product-service:8002
      ORDER_SERVICE_URL: http://order-service:8003
    depends_on:
      - user-service
      - product-service
      - order-service

volumes:
  postgres_data:
```

## Running the Example

### With Docker Compose

```bash
# Start all services
docker-compose up -d

# Check service health
curl http://localhost:8000/health

# View logs
docker-compose logs -f
```

### Locally

```bash
# Terminal 1 - User Service
cd user-service
go run main.go

# Terminal 2 - Product Service
cd product-service
go run main.go

# Terminal 3 - Order Service
cd order-service
go run main.go

# Terminal 4 - API Gateway
cd api-gateway
go run main.go
```

## API Examples

### Create User

```bash
curl -X POST http://localhost:8000/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

### Create Product

```bash
curl -X POST http://localhost:8000/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro",
    "description": "16-inch laptop",
    "price": 2499.99,
    "stock": 10
  }'
```

### Create Order

```bash
curl -X POST http://localhost:8000/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-uuid-here",
    "items": [
      {
        "product_id": "product-uuid-here",
        "quantity": 1
      }
    ]
  }'
```

## Best Practices

### 1. Service Independence

Each service should:

- Have its own database
- Be deployable independently
- Have its own configuration

### 2. Error Handling

Handle network failures gracefully:

```go
func (c *ProductClient) GetProduct(id string) (*Product, error) {
    retries := 3
    for i := 0; i < retries; i++ {
        product, err := c.fetchProduct(id)
        if err == nil {
            return product, nil
        }

        if i < retries-1 {
            time.Sleep(time.Second * time.Duration(i+1))
        }
    }
    return nil, errors.New("max retries exceeded")
}
```

### 3. Circuit Breaker

Prevent cascading failures:

```go
import "github.com/sony/gobreaker"

cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:        "ProductService",
    MaxRequests: 3,
    Timeout:     30 * time.Second,
})

product, err := cb.Execute(func() (interface{}, error) {
    return c.GetProduct(id)
})
```

## Next Steps

- Add [Service Discovery with Consul](./service-discovery)
- Implement [Event-Driven Architecture](./event-driven)
- Add [Distributed Tracing](./tracing)
- Learn about [API Gateway Patterns](./api-gateway)

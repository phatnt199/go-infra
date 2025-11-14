---
sidebar_position: 2
---

# Building CRUD Operations

Learn how to build complete REST APIs with CRUD operations using go-infra's components.

> **Note**: go-infra does not provide built-in automatic CRUD generation. Instead, it provides powerful building blocks that make implementing CRUD operations straightforward and maintainable.

## Overview

This guide shows you how to build CRUD (Create, Read, Update, Delete) operations using:

- **Fiber** for HTTP routing
- **GORM** for database operations
- **Validator** for request validation
- **Pagination utilities** for list endpoints
- **Mapper** for DTO transformations

## Complete Example

### 1. Define Your Domain Model

```go
package domain

import "github.com/phatnt199/go-infra/pkg/domain/entity"

type Product struct {
    entity.BaseModel
    Name        string  `json:"name" gorm:"not null"`
    Description string  `json:"description"`
    Price       float64 `json:"price" gorm:"not null"`
    Stock       int     `json:"stock" gorm:"default:0"`
}
```

### 2. Create DTOs

```go
package dto

type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required"`
    Description string  `json:"description"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Stock       int     `json:"stock" validate:"gte=0"`
}

type UpdateProductRequest struct {
    Name        *string  `json:"name,omitempty"`
    Description *string  `json:"description,omitempty"`
    Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
    Stock       *int     `json:"stock,omitempty" validate:"omitempty,gte=0"`
}

type ProductResponse struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
    CreatedAt   string  `json:"created_at"`
    UpdatedAt   string  `json:"updated_at"`
}
```

### 3. Create Handler

```go
package handler

import (
    "github.com/gofiber/fiber/v2"
    "github.com/phatnt199/go-infra/pkg/utils"
    "github.com/phatnt199/go-infra/pkg/validator"
    "github.com/phatnt199/go-infra/pkg/mapper"
    "gorm.io/gorm"
    "myapp/internal/domain"
    "myapp/internal/dto"
)

type ProductHandler struct {
    db *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
    return &ProductHandler{db: db}
}

// Create - POST /api/products
func (h *ProductHandler) Create(c *fiber.Ctx) error {
    var req dto.CreateProductRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Validate request
    if err := validator.Struct(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    // Map DTO to domain model
    product := &domain.Product{
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
    }

    // Save to database
    if err := h.db.Create(product).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to create product",
        })
    }

    // Map to response
    response, _ := mapper.Map[dto.ProductResponse](product)
    return c.Status(fiber.StatusCreated).JSON(response)
}

// GetByID - GET /api/products/:id
func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
    id := c.Params("id")

    var product domain.Product
    if err := h.db.First(&product, "id = ?", id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "Product not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch product",
        })
    }

    response, _ := mapper.Map[dto.ProductResponse](product)
    return c.JSON(response)
}

// List - GET /api/products
func (h *ProductHandler) List(c *fiber.Ctx) error {
    // Parse pagination parameters
    listQuery, err := utils.GetListQueryFromContext(c)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid query parameters",
        })
    }

    var products []domain.Product
    var total int64

    // Count total records
    h.db.Model(&domain.Product{}).Count(&total)

    // Fetch paginated results
    if err := h.db.
        Limit(listQuery.GetLimit()).
        Offset(listQuery.GetOffset()).
        Order(listQuery.GetOrderBy()).
        Find(&products).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch products",
        })
    }

    // Map to response DTOs
    responses, _ := mapper.Map[[]dto.ProductResponse](products)

    // Create list result with pagination
    result := utils.NewListResult(
        responses,
        listQuery.GetSize(),
        listQuery.GetPage(),
        total,
    )

    return c.JSON(result)
}

// Update - PUT /api/products/:id
func (h *ProductHandler) Update(c *fiber.Ctx) error {
    id := c.Params("id")

    var req dto.UpdateProductRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Validate request
    if err := validator.Struct(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    // Find existing product
    var product domain.Product
    if err := h.db.First(&product, "id = ?", id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "Product not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch product",
        })
    }

    // Update fields
    if req.Name != nil {
        product.Name = *req.Name
    }
    if req.Description != nil {
        product.Description = *req.Description
    }
    if req.Price != nil {
        product.Price = *req.Price
    }
    if req.Stock != nil {
        product.Stock = *req.Stock
    }

    // Save changes
    if err := h.db.Save(&product).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to update product",
        })
    }

    response, _ := mapper.Map[dto.ProductResponse](product)
    return c.JSON(response)
}

// Delete - DELETE /api/products/:id
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
    id := c.Params("id")

    result := h.db.Delete(&domain.Product{}, "id = ?", id)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to delete product",
        })
    }

    if result.RowsAffected == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Product not found",
        })
    }

    return c.Status(fiber.StatusNoContent).Send(nil)
}
```

### 4. Register Routes

```go
package main

import (
    "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
    "github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
    "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
    "myapp/internal/handler"
    "go.uber.org/fx"
)

func main() {
    app := fxapp.NewApplicationBuilder().
        ProvideModule(fiber_adapter.Module).
        ProvideModule(gorm.Module).
        Provide(handler.NewProductHandler).
        Provide(fx.Invoke(registerRoutes)).
        Build()

    app.Run()
}

func registerRoutes(
    server contracts.HttpServer,
    productHandler *handler.ProductHandler,
) {
    server.RouteBuilder().RegisterHandler(func(router interface{}) {
        if app, ok := router.(*fiber.App); ok {
            api := app.Group("/api")
            products := api.Group("/products")

            products.Post("/", productHandler.Create)
            products.Get("/:id", productHandler.GetByID)
            products.Get("/", productHandler.List)
            products.Put("/:id", productHandler.Update)
            products.Delete("/:id", productHandler.Delete)
        }
    })
}
```

## API Endpoints

### Create Product

**Request:**

```bash
POST /api/products
Content-Type: application/json

{
  "name": "Laptop",
  "description": "Gaming laptop",
  "price": 1299.99,
  "stock": 10
}
```

**Response:** `201 Created`

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440000",
	"name": "Laptop",
	"description": "Gaming laptop",
	"price": 1299.99,
	"stock": 10,
	"created_at": "2024-01-01T00:00:00Z",
	"updated_at": "2024-01-01T00:00:00Z"
}
```

### Get Product by ID

**Request:**

```bash
GET /api/products/550e8400-e29b-41d4-a716-446655440000
```

**Response:** `200 OK`

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440000",
	"name": "Laptop",
	"description": "Gaming laptop",
	"price": 1299.99,
	"stock": 10,
	"created_at": "2024-01-01T00:00:00Z",
	"updated_at": "2024-01-01T00:00:00Z"
}
```

### List Products (Paginated)

**Request:**

```bash
GET /api/products?page=1&size=10&orderBy=created_at DESC
```

**Response:** `200 OK`

```json
{
	"items": [
		{
			"id": "550e8400-e29b-41d4-a716-446655440000",
			"name": "Laptop",
			"price": 1299.99,
			"stock": 10
		}
	],
	"page": 1,
	"size": 10,
	"totalItems": 45,
	"totalPage": 5
}
```

### Update Product

**Request:**

```bash
PUT /api/products/550e8400-e29b-41d4-a716-446655440000
Content-Type: application/json

{
  "price": 1199.99,
  "stock": 5
}
```

**Response:** `200 OK`

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440000",
	"name": "Laptop",
	"description": "Gaming laptop",
	"price": 1199.99,
	"stock": 5,
	"created_at": "2024-01-01T00:00:00Z",
	"updated_at": "2024-01-01T12:00:00Z"
}
```

### Delete Product

**Request:**

```bash
DELETE /api/products/550e8400-e29b-41d4-a716-446655440000
```

**Response:** `204 No Content`

## Best Practices

### 1. Use DTOs for API Contracts

Separate your domain models from API requests/responses:

```go
// Domain model - internal representation
type Product struct {
    entity.BaseModel
    Name  string
    Price float64
}

// API DTO - external contract
type ProductDTO struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

### 2. Validate All Inputs

Always validate request data:

```go
import "github.com/phatnt199/go-infra/pkg/validator"

if err := validator.Struct(&req); err != nil {
    return c.Status(400).JSON(fiber.Map{"error": err.Error()})
}
```

### 3. Use Mapper for Transformations

Convert between types safely:

```go
import "github.com/phatnt199/go-infra/pkg/mapper"

// Configure mapping once
mapper.CreateMap[domain.Product, dto.ProductResponse]()

// Use everywhere
response, err := mapper.Map[dto.ProductResponse](product)
```

### 4. Implement Pagination

Use built-in utilities:

```go
import "github.com/phatnt199/go-infra/pkg/utils"

listQuery, _ := utils.GetListQueryFromContext(c)
result := utils.NewListResult(items, listQuery.GetSize(), listQuery.GetPage(), total)
```

### 5. Handle Errors Consistently

Create standard error responses:

```go
func handleError(c *fiber.Ctx, status int, message string) error {
    return c.Status(status).JSON(fiber.Map{
        "error": message,
        "timestamp": time.Now().Unix(),
    })
}
```

## Next Steps

- **[Middleware](./middleware)** - Add authentication, logging
- **[Database Migrations](../database/migrations)** - Manage schema changes
- **[Authentication](../authentication/getting-started)** - Secure your endpoints

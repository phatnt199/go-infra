---
sidebar_position: 1
---

# Database Setup

Learn how to integrate PostgreSQL with go-infra using GORM.

## Quick Start

### 1. Add Database Module

```go
import (
    "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
    "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
)

func main() {
    app := fxapp.NewApplicationBuilder().
        ProvideModule(gorm.Module).
        Build()

    app.Run()
}
```

### 2. Configuration

Set database credentials in `.env`:

```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=myapp
POSTGRES_SSLMODE=disable
POSTGRES_MAX_OPEN_CONNS=25
POSTGRES_MAX_IDLE_CONNS=5
POSTGRES_CONN_MAX_LIFETIME=5m
```

## Define Models

### Basic Model

```go
package domain

import "github.com/phatnt199/go-infra/pkg/domain/entity"

type User struct {
    entity.BaseModel          // ID, CreatedAt, UpdatedAt, DeletedAt
    Name  string `json:"name" gorm:"not null"`
    Email string `json:"email" gorm:"uniqueIndex;not null"`
    Age   int    `json:"age"`
}
```

The `BaseModel` provides:

- `ID` - UUID primary key
- `CreatedAt` - Timestamp
- `UpdatedAt` - Timestamp
- `DeletedAt` - For soft deletes

### Custom Model

```go
type Product struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    Price       float64   `json:"price" gorm:"not null"`
    Stock       int       `json:"stock" gorm:"default:0"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Relationships

```go
// One-to-Many
type User struct {
    entity.BaseModel
    Name   string    `json:"name"`
    Posts  []Post    `json:"posts" gorm:"foreignKey:UserID"`
}

type Post struct {
    entity.BaseModel
    Title  string `json:"title"`
    Body   string `json:"body"`
    UserID string `json:"user_id"`
    User   User   `json:"user" gorm:"foreignKey:UserID"`
}

// Many-to-Many
type User struct {
    entity.BaseModel
    Name  string  `json:"name"`
    Roles []Role  `json:"roles" gorm:"many2many:user_roles"`
}

type Role struct {
    entity.BaseModel
    Name  string  `json:"name"`
    Users []User  `json:"users" gorm:"many2many:user_roles"`
}
```

## Auto Migration

### Simple Migration

```go
func main() {
    app := fxapp.NewApplicationBuilder().
        ProvideModule(gorm.Module).
        Provide(fx.Invoke(runMigrations)).
        Build()

    app.Run()
}

func runMigrations(db *gorm.DB) {
    db.AutoMigrate(
        &domain.User{},
        &domain.Post{},
        &domain.Product{},
    )
}
```

### With Logging

```go
func runMigrations(db *gorm.DB, logger logger.Logger) {
    logger.Info("Running migrations...")

    if err := db.AutoMigrate(
        &domain.User{},
        &domain.Post{},
    ); err != nil {
        logger.Fatal("Migration failed", logger.Err(err))
    }

    logger.Info("Migrations completed")
}
```

## CRUD Operations

### Create

```go
func CreateUser(db *gorm.DB, name, email string) (*domain.User, error) {
    user := &domain.User{
        Name:  name,
        Email: email,
    }

    result := db.Create(user)
    if result.Error != nil {
        return nil, result.Error
    }

    return user, nil
}
```

### Read

```go
// Find by ID
func GetUser(db *gorm.DB, id string) (*domain.User, error) {
    var user domain.User
    result := db.First(&user, "id = ?", id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

// Find by email
func GetUserByEmail(db *gorm.DB, email string) (*domain.User, error) {
    var user domain.User
    result := db.Where("email = ?", email).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

// Find all
func GetAllUsers(db *gorm.DB) ([]domain.User, error) {
    var users []domain.User
    result := db.Find(&users)
    if result.Error != nil {
        return nil, result.Error
    }
    return users, nil
}

// With conditions
func GetActiveUsers(db *gorm.DB) ([]domain.User, error) {
    var users []domain.User
    result := db.Where("status = ?", "active").Find(&users)
    return users, result.Error
}
```

### Update

```go
// Update specific fields
func UpdateUser(db *gorm.DB, id string, updates map[string]interface{}) error {
    result := db.Model(&domain.User{}).
        Where("id = ?", id).
        Updates(updates)
    return result.Error
}

// Update struct
func UpdateUserStruct(db *gorm.DB, user *domain.User) error {
    result := db.Save(user)
    return result.Error
}
```

### Delete

```go
// Soft delete (with DeletedAt)
func DeleteUser(db *gorm.DB, id string) error {
    result := db.Delete(&domain.User{}, "id = ?", id)
    return result.Error
}

// Hard delete (permanent)
func HardDeleteUser(db *gorm.DB, id string) error {
    result := db.Unscoped().Delete(&domain.User{}, "id = ?", id)
    return result.Error
}
```

## Repository Pattern

### Define Repository

```go
package repository

import (
    "context"
    "gorm.io/gorm"
    "myapp/internal/domain"
)

type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id string) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    FindAll(ctx context.Context) ([]domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id string) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]domain.User, error) {
    var users []domain.User
    err := r.db.WithContext(ctx).Find(&users).Error
    return users, err
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error
}
```

### Use Repository

```go
func main() {
    app := fxapp.NewApplicationBuilder().
        ProvideModule(gorm.Module).
        Provide(repository.NewUserRepository).
        Provide(service.NewUserService).
        Build()

    app.Run()
}

// Service using repository
type UserService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
    user := &domain.User{
        Name:  name,
        Email: email,
    }

    err := s.repo.Create(ctx, user)
    if err != nil {
        return nil, err
    }

    return user, nil
}
```

## Queries

### Basic Queries

```go
// Simple where
db.Where("name = ?", "John").Find(&users)

// Multiple conditions
db.Where("name = ? AND age > ?", "John", 18).Find(&users)

// IN clause
db.Where("id IN ?", []string{"1", "2", "3"}).Find(&users)

// LIKE
db.Where("email LIKE ?", "%@example.com").Find(&users)

// Order by
db.Order("created_at desc").Find(&users)

// Limit and offset
db.Limit(10).Offset(20).Find(&users)
```

### Complex Queries

```go
// Group by
var results []struct {
    Status string
    Count  int
}
db.Model(&domain.User{}).
    Select("status, count(*) as count").
    Group("status").
    Scan(&results)

// Joins
var users []domain.User
db.Joins("LEFT JOIN posts ON posts.user_id = users.id").
    Where("posts.status = ?", "published").
    Find(&users)

// Preload relationships
db.Preload("Posts").Find(&users)
db.Preload("Posts.Comments").Find(&users)
```

### Raw SQL

```go
// Raw query
var users []domain.User
db.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)

// Named parameters
db.Raw("SELECT * FROM users WHERE name = @name AND age > @age",
    sql.Named("name", "John"),
    sql.Named("age", 18),
).Scan(&users)

// Exec
db.Exec("UPDATE users SET status = ? WHERE id = ?", "active", id)
```

## Transactions

### Basic Transaction

```go
func CreateUserWithProfile(db *gorm.DB, user *domain.User, profile *domain.Profile) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // Create user
        if err := tx.Create(user).Error; err != nil {
            return err
        }

        // Create profile
        profile.UserID = user.ID
        if err := tx.Create(profile).Error; err != nil {
            return err
        }

        return nil
    })
}
```

### Manual Transaction

```go
func TransferMoney(db *gorm.DB, fromID, toID string, amount float64) error {
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Deduct from sender
    if err := tx.Model(&domain.Account{}).
        Where("id = ?", fromID).
        Update("balance", gorm.Expr("balance - ?", amount)).
        Error; err != nil {
        tx.Rollback()
        return err
    }

    // Add to receiver
    if err := tx.Model(&domain.Account{}).
        Where("id = ?", toID).
        Update("balance", gorm.Expr("balance + ?", amount)).
        Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}
```

## Pagination

```go
type PaginationParams struct {
    Page  int
    Limit int
}

type PaginatedResult struct {
    Data       interface{} `json:"data"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    Limit      int         `json:"limit"`
    TotalPages int         `json:"total_pages"`
}

func PaginateUsers(db *gorm.DB, params PaginationParams) (*PaginatedResult, error) {
    var users []domain.User
    var total int64

    // Count total
    if err := db.Model(&domain.User{}).Count(&total).Error; err != nil {
        return nil, err
    }

    // Get page data
    offset := (params.Page - 1) * params.Limit
    if err := db.Limit(params.Limit).Offset(offset).Find(&users).Error; err != nil {
        return nil, err
    }

    totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

    return &PaginatedResult{
        Data:       users,
        Total:      total,
        Page:       params.Page,
        Limit:      params.Limit,
        TotalPages: totalPages,
    }, nil
}
```

## Error Handling

```go
import (
    "errors"
    "gorm.io/gorm"
)

func GetUser(db *gorm.DB, id string) (*domain.User, error) {
    var user domain.User
    err := db.First(&user, "id = ?", id).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("user not found")
        }
        return nil, errors.Wrap(err, "database error")
    }

    return &user, nil
}
```

## Best Practices

### 1. Always Use Context

```go
func (r *Repository) FindByID(ctx context.Context, id string) (*User, error) {
    return r.db.WithContext(ctx).First(&User{}, "id = ?", id)
}
```

### 2. Use Transactions for Multiple Operations

```go
db.Transaction(func(tx *gorm.DB) error {
    // Multiple operations here
    return nil
})
```

### 3. Handle Errors Properly

```go
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        // Handle not found
    }
    return err
}
```

### 4. Use Preloading Wisely

```go
// Good: Load only what you need
db.Preload("Posts").Find(&users)

// Bad: N+1 query problem
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&user.Posts)
}
```

### 5. Index Your Queries

```go
type User struct {
    entity.BaseModel
    Email string `gorm:"uniqueIndex"`
    Name  string `gorm:"index"`
}
```

## Next Steps

- **[Migrations](./migrations)** - Database migration strategies
- **[Advanced Queries](./advanced-queries)** - Complex query patterns
- **[Performance](./performance)** - Optimize database access

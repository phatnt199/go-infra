---
sidebar_position: 5
---

# Utilities

go-infra provides a rich set of utility packages to simplify common development tasks.

## Pagination

The `utils` package provides built-in pagination support for list endpoints.

### ListQuery

Parse and manage pagination parameters from HTTP requests:

```go
import "github.com/phatnt199/go-infra/pkg/utils"

func (h *Handler) List(c *fiber.Ctx) error {
    // Parse pagination from query params: ?page=1&size=10&orderBy=created_at DESC
    listQuery, err := utils.GetListQueryFromContext(c)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid query parameters"})
    }

    // Use in database query
    var items []MyModel
    db.Limit(listQuery.GetLimit()).
       Offset(listQuery.GetOffset()).
       Order(listQuery.GetOrderBy()).
       Find(&items)

    // Get pagination info
    page := listQuery.GetPage()      // Current page number
    size := listQuery.GetSize()      // Page size
    offset := listQuery.GetOffset()  // Offset for SQL query
    limit := listQuery.GetLimit()    // Limit for SQL query
}
```

### ListResult

Create paginated responses:

```go
var items []Product
var total int64

db.Model(&Product{}).Count(&total)
db.Limit(size).Offset(offset).Find(&items)

// Create paginated result
result := utils.NewListResult(
    items,           // Items for current page
    size,            // Page size
    page,            // Current page number
    total,           // Total item count
)

return c.JSON(result)
```

**Response format:**

```json
{
  "items": [...],
  "page": 1,
  "size": 10,
  "totalItems": 45,
  "totalPage": 5
}
```

### Converting with Mapper

Transform paginated results to DTOs:

```go
import "github.com/phatnt199/go-infra/pkg/mapper"

// Get domain models with pagination
domainResult := utils.NewListResult(products, size, page, total)

// Convert to DTOs
dtoResult, err := utils.ListResultToListResultDto[dto.ProductResponse](domainResult)
```

## Validator

Wrapper around `go-playground/validator` for struct validation.

### Basic Usage

```go
import "github.com/phatnt199/go-infra/pkg/validator"

type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"required,gte=18,lte=120"`
    Password string `json:"password" validate:"required,min=8"`
}

func CreateUser(c *fiber.Ctx) error {
    var req CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    // Validate struct
    if err := validator.Struct(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }

    // Validation passed, proceed...
}
```

### Validation Tags

Common validation tags:

```go
type Example struct {
    // Required field
    Name string `validate:"required"`

    // String length
    Username string `validate:"min=3,max=20"`

    // Numeric range
    Age int `validate:"gte=0,lte=120"`

    // Email format
    Email string `validate:"email"`

    // URL format
    Website string `validate:"url"`

    // Enum values
    Status string `validate:"oneof=active inactive pending"`

    // Greater than field
    EndDate time.Time `validate:"gtefield=StartDate"`

    // Optional field with validation
    Phone string `validate:"omitempty,e164"`
}
```

### Custom Validators

```go
import "github.com/go-playground/validator/v10"

// Register custom validator
v := validator.New()
v.RegisterValidation("custom", func(fl validator.FieldLevel) bool {
    return fl.Field().String() != "forbidden"
})
```

## Mapper

Type-safe object-to-object mapping with generics.

### Setup

```go
import "github.com/phatnt199/go-infra/pkg/mapper"

// Configure mappings at application startup
func init() {
    mapper.CreateMap[domain.User, dto.UserResponse]()
    mapper.CreateMap[domain.Product, dto.ProductResponse]()
}
```

### Basic Mapping

```go
// Map single object
user := &domain.User{Name: "John", Email: "john@example.com"}
response, err := mapper.Map[dto.UserResponse](user)

// Map slice
users := []domain.User{...}
responses, err := mapper.Map[[]dto.UserResponse](users)
```

### Custom Mapping

For complex transformations, use custom map functions:

```go
mapper.CreateCustomMap(func(src domain.User) dto.UserResponse {
    return dto.UserResponse{
        ID:       src.ID,
        FullName: src.FirstName + " " + src.LastName,
        Email:    strings.ToLower(src.Email),
        Age:      calculateAge(src.BirthDate),
    }
})
```

### Field Mapping with Tags

Use `mapper` tags to map fields with different names:

```go
type Source struct {
    UserName string `mapper:"user_name"`
    UserAge  int    `mapper:"age"`
}

type Destination struct {
    UserName string `json:"user_name"`
    Age      int    `json:"age"`
}
```

### Configuration

```go
import "github.com/phatnt199/go-infra/pkg/mapper"

// Configure mapper behavior
mapper.Configure(&mapper.MapperConfig{
    MapUnexportedFields: false, // Don't map private fields
})

// Clear all mappings (useful in tests)
mapper.ClearMappings()
```

## Crypto

Security utilities for JWT, password hashing, and encryption.

### JWT Manager

```go
import "github.com/phatnt199/go-infra/pkg/crypto"

// Create JWT manager
config := &crypto.JWTConfig{
    Secret:             "your-secret-key",
    Algorithm:          crypto.AlgorithmHS256,
    Issuer:             "myapp",
    Audience:           "myapp-api",
    AccessTokenExpiry:  15 * time.Minute,
    RefreshTokenExpiry: 7 * 24 * time.Hour,
}

jwtManager, err := crypto.NewJWTManager(config)

// Generate tokens
claims := &crypto.Claims{
    UserID:   "user-123",
    Username: "john",
    Email:    "john@example.com",
    Roles:    []string{"user", "admin"},
}

accessToken, err := jwtManager.GenerateToken(claims, crypto.AccessToken)
refreshToken, err := jwtManager.GenerateToken(claims, crypto.RefreshToken)

// Or generate both at once
accessToken, refreshToken, err := jwtManager.GenerateTokenPair(claims)

// Verify and parse token
claims, err := jwtManager.ParseToken(tokenString)
if err != nil {
    // Token is invalid or expired
}

// Simple validation
err = jwtManager.ValidateToken(tokenString)
```

### RSA Keys for JWT

For production, use RSA instead of HMAC:

```go
import "github.com/phatnt199/go-infra/pkg/crypto"

// Load keys from files
privateKey, err := crypto.LoadRSAPrivateKeyFromFile("./keys/private.pem")
publicKey, err := crypto.LoadRSAPublicKeyFromFile("./keys/public.pem")

config := &crypto.JWTConfig{
    Algorithm:          crypto.AlgorithmRS256,
    PrivateKey:         privateKey,
    PublicKey:          publicKey,
    Issuer:             "myapp",
    Audience:           "myapp-api",
    AccessTokenExpiry:  15 * time.Minute,
    RefreshTokenExpiry: 7 * 24 * time.Hour,
}

jwtManager, err := crypto.NewJWTManager(config)
```

### Password Hashing

```go
import "github.com/phatnt199/go-infra/pkg/crypto"

// Hash password with bcrypt
hashedPassword, err := crypto.HashPasswordBcrypt("user-password")

// Hash with Argon2 (more secure)
hashedPassword, err := crypto.HashPasswordArgon2("user-password")

// Verify password
isValid, err := crypto.ComparePassword("user-password", hashedPassword)
if !isValid {
    return errors.New("invalid password")
}

// Custom hasher with configuration
config := &crypto.HashConfig{
    BcryptCost:    12,  // Higher = more secure, slower
    Argon2Time:    1,
    Argon2Memory:  64 * 1024,
    Argon2Threads: 4,
    Argon2KeyLen:  32,
}

hasher := crypto.NewHasher(config)
hash, err := hasher.HashPassword("password", crypto.AlgorithmArgon2)
```

## Reflection Helpers

Advanced reflection utilities for working with structs.

### Get Field Values

```go
import reflectionHelper "github.com/phatnt199/go-infra/pkg/reflection/reflection_helper"

type User struct {
    Name  string
    email string // unexported
}

user := User{Name: "John", email: "john@example.com"}

// Get by field name (works with unexported fields)
name := reflectionHelper.GetFieldValueByName(user, "Name")
email := reflectionHelper.GetFieldValueByName(user, "email")

// Get by index
firstField := reflectionHelper.GetFieldValueByIndex(user, 0)
```

### Set Field Values

```go
// Set exported fields
reflectionHelper.SetFieldValueByName(&user, "Name", "Jane")

// Set unexported fields
reflectionHelper.SetFieldValueByName(&user, "email", "jane@example.com")
```

### Get All Fields

```go
fields := reflectionHelper.GetAllFields(reflect.TypeOf(User{}))
for _, field := range fields {
    fmt.Printf("Field: %s, Type: %s\n", field.Name, field.Type)
}
```

### Type Information

```go
// Get type path
path := reflectionHelper.ObjectTypePath(&user)
// Output: "myapp/domain.User"

// Get type path from generic
path := reflectionHelper.TypePath[User]()

// Get method path
path := reflectionHelper.MethodPath(someFunction)
```

## Common Utilities

### Pointer Helpers

```go
import "github.com/phatnt199/go-infra/pkg/utils"

// Create pointers to values
name := utils.ToPtr("John")
age := utils.ToPtr(30)
isActive := utils.ToPtr(true)

// Useful for optional fields in update requests
type UpdateRequest struct {
    Name  *string `json:"name,omitempty"`
    Age   *int    `json:"age,omitempty"`
}
```

### String Utilities

```go
import "github.com/phatnt199/go-infra/pkg/utils"

// Check if string is empty
if utils.IsEmpty(str) {
    // ...
}

// Generate random string
randomStr := utils.RandomString(16)

// Truncate string
short := utils.Truncate(longStr, 100)
```

### Slice Utilities

```go
import "github.com/phatnt199/go-infra/pkg/utils"

// Contains
if utils.Contains([]string{"a", "b", "c"}, "b") {
    // ...
}

// Remove duplicates
unique := utils.RemoveDuplicates([]string{"a", "b", "a", "c"})

// Filter
filtered := utils.Filter(items, func(item Item) bool {
    return item.Active
})
```

### Error Utilities

```go
import "github.com/phatnt199/go-infra/pkg/utils"

// Wrap errors with context
err := utils.WrapError(err, "failed to create user")

// Chain multiple errors
errs := utils.NewErrorChain()
errs.Add(err1)
errs.Add(err2)
if errs.HasErrors() {
    return errs.Error()
}
```

## Next Steps

- **[Configuration](./configuration)** - Environment-based config
- **[Health Checks](./health-checks)** - Service health monitoring
- **[Logging](./logging)** - Structured logging
- **[Testing](../advanced/testing)** - Testing utilities

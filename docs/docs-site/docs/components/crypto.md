---
sidebar_position: 2
---

# Crypto

Production-ready cryptographic utilities for password hashing, JWT tokens, and encryption.

## Features

- ✅ **Password hashing** - Bcrypt and Argon2
- ✅ **JWT tokens** - HMAC and RSA support
- ✅ **AES encryption** - AES-128/192/256-GCM
- ✅ **Easy to use** - Simple, secure defaults
- ✅ **Battle-tested** - Production-proven

## Password Hashing

### Bcrypt (Battle-tested)

```go
import "github.com/phatnt199/go-infra/pkg/crypto"

// Hash password
hash, err := crypto.HashPasswordBcrypt("MyPassword123")

// Verify password
match, err := crypto.ComparePassword("MyPassword123", hash)
if match {
    // Password is correct
}
```

### Argon2 (Modern, recommended)

```go
// Hash password (OWASP recommended)
hash, err := crypto.HashPasswordArgon2("MyPassword123")

// Verify password (auto-detects algorithm)
match, err := crypto.ComparePassword("MyPassword123", hash)
```

### Custom Configuration

```go
hasher := crypto.NewHasher(&crypto.HashConfig{
    BcryptCost:    12,        // Higher = more secure, slower
    Argon2Time:    2,         // Number of iterations
    Argon2Memory:  64 * 1024, // 64 MB
    Argon2Threads: 4,         // Parallelism
    Argon2KeyLen:  32,        // Output length
})

hash, err := hasher.HashPassword("password", crypto.AlgorithmArgon2)
```

## JWT Tokens

### Setup

```go
import "github.com/phatnt199/go-infra/pkg/crypto"

// HMAC configuration (symmetric key)
config := &crypto.JWTConfig{
    Secret:             "your-secret-key-here",
    Algorithm:          crypto.AlgorithmHS256,
    Issuer:             "your-app",
    Audience:           "your-api",
    AccessTokenExpiry:  15 * time.Minute,
    RefreshTokenExpiry: 7 * 24 * time.Hour,
}

crypto.InitDefaultJWT(config)
```

### Generate Tokens

```go
// Create claims
claims := &crypto.Claims{
    UserID:   "user123",
    Username: "johndoe",
    Email:    "john@example.com",
    Roles:    []string{"admin", "user"},
    Custom: map[string]interface{}{
        "department": "engineering",
    },
}

// Generate access token
accessToken, err := crypto.GenerateAccessToken(claims)

// Generate refresh token
refreshToken, err := crypto.GenerateRefreshToken(claims)

// Generate both at once
accessToken, refreshToken, err := manager.GenerateTokenPair(claims)
```

### Validate Tokens

```go
// Parse and validate token
claims, err := crypto.ParseJWT(tokenString)
if err != nil {
    // Handle invalid token
    if errors.Is(err, errors.CodeTokenExpired) {
        // Token expired
    }
}

// Access claims
userID := claims.UserID
roles := claims.Roles
```

### Refresh Tokens

```go
// Refresh an access token using a refresh token
newAccessToken, err := manager.RefreshToken(refreshToken)
```

## Encryption

### Generate Keys

```go
// Generate AES-256 key (recommended)
key, err := crypto.GenerateAES256Key()

// Convert key to string for storage
keyStr := crypto.KeyToString(key)

// Convert string back to key
key, err := crypto.KeyFromString(keyStr)
```

### Encrypt/Decrypt Strings

```go
// Initialize encryptor
key, _ := crypto.GenerateAES256Key()
crypto.InitDefaultEncryptor(key)

// Encrypt string
plaintext := "Sensitive data"
ciphertext, err := crypto.Encrypt(plaintext)

// Decrypt string
decrypted, err := crypto.Decrypt(ciphertext)
```

### Encrypt/Decrypt Bytes

```go
// Encrypt bytes
data := []byte("Binary data")
encrypted, err := crypto.EncryptBytes(data)

// Decrypt bytes
decrypted, err := crypto.DecryptBytes(encrypted)
```

## Complete Auth Example

### User Registration

```go
func RegisterUser(username, password string) error {
    // Validate password strength
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }

    // Hash password
    hash, err := crypto.HashPasswordArgon2(password)
    if err != nil {
        return errors.Wrap(err, "failed to hash password")
    }

    // Store user
    user := &User{
        Username: username,
        Password: hash,
    }

    return db.Create(user).Error
}
```

### User Login

```go
func LoginUser(username, password string) (string, error) {
    // Fetch user
    var user User
    if err := db.Where("username = ?", username).First(&user).Error; err != nil {
        return "", errors.New("invalid credentials")
    }

    // Verify password
    match, err := crypto.ComparePassword(password, user.Password)
    if err != nil || !match {
        return "", errors.New("invalid credentials")
    }

    // Generate JWT token
    claims := &crypto.Claims{
        UserID:   user.ID,
        Username: user.Username,
        Roles:    user.Roles,
    }

    token, err := crypto.GenerateAccessToken(claims)
    if err != nil {
        return "", errors.Wrap(err, "failed to generate token")
    }

    return token, nil
}
```

### Auth Middleware

```go
func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Extract token
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "missing authorization header",
            })
        }

        // Format: "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Status(401).JSON(fiber.Map{
                "error": "invalid authorization format",
            })
        }

        token := parts[1]

        // Validate token
        claims, err := crypto.ParseJWT(token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "invalid token",
            })
        }

        // Add claims to context
        c.Locals("user_id", claims.UserID)
        c.Locals("username", claims.Username)
        c.Locals("roles", claims.Roles)

        return c.Next()
    }
}
```

## Encrypting Database Fields

```go
type User struct {
    ID       uint
    Username string
    Email    string
    SSN      string // Encrypted field
}

func (u *User) EncryptSSN(ssn string) error {
    encrypted, err := crypto.Encrypt(ssn)
    if err != nil {
        return err
    }
    u.SSN = encrypted
    return nil
}

func (u *User) DecryptSSN() (string, error) {
    return crypto.Decrypt(u.SSN)
}

// Usage
user := &User{Username: "john", Email: "john@example.com"}

// Encrypt before saving
user.EncryptSSN("123-45-6789")
db.Create(user)

// Decrypt when reading
ssn, err := user.DecryptSSN()
```

## RSA Keys for JWT

### Generate RSA Keys

```bash
# Generate private key
openssl genrsa -out private.pem 2048

# Generate public key
openssl rsa -in private.pem -pubout -out public.pem
```

### Use RSA in JWT

```go
// Load keys
privateKey, err := crypto.LoadRSAPrivateKeyFromFile("private.pem")
publicKey, err := crypto.LoadRSAPublicKeyFromFile("public.pem")

// Configure JWT with RSA
config := &crypto.JWTConfig{
    PrivateKey:         privateKey,
    PublicKey:          publicKey,
    Algorithm:          crypto.AlgorithmRS256,
    Issuer:             "your-app",
    AccessTokenExpiry:  15 * time.Minute,
}

crypto.InitDefaultJWT(config)
```

## Security Best Practices

### 1. Strong Secrets

```bash
# Generate secure secret (32+ bytes)
openssl rand -base64 32
```

### 2. Password Requirements

```go
func ValidatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }

    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

    if !hasUpper || !hasLower || !hasNumber {
        return errors.New("password must contain uppercase, lowercase, and number")
    }

    return nil
}
```

### 3. Short Token Expiry

```go
AccessTokenExpiry:  15 * time.Minute  // Short-lived
RefreshTokenExpiry: 7 * 24 * time.Hour // Longer-lived
```

### 4. Store Keys Securely

```bash
# Environment variables
export JWT_SECRET=$(openssl rand -base64 32)
export ENCRYPTION_KEY=$(openssl rand -base64 32)

# Or use key management services
# - AWS KMS
# - HashiCorp Vault
# - Google Cloud KMS
```

### 5. Use HTTPS

Always use HTTPS in production. Never send tokens over HTTP.

### 6. Rotate Keys

Plan for periodic key rotation:

- JWT secrets: Every 6-12 months
- Encryption keys: Every year
- Implement key versioning

## Algorithm Comparison

### Password Hashing

| Algorithm | Speed  | Memory | Security | Recommendation    |
| --------- | ------ | ------ | -------- | ----------------- |
| Bcrypt    | Slow   | Low    | Good     | Production-ready  |
| Argon2    | Slower | High   | Best     | OWASP recommended |

**Use Argon2 for new projects**, Bcrypt is still secure.

### JWT Algorithms

| Algorithm | Type | Speed  | Use Case                         |
| --------- | ---- | ------ | -------------------------------- |
| HS256     | HMAC | Fast   | Single service                   |
| HS512     | HMAC | Fast   | Single service, higher security  |
| RS256     | RSA  | Slower | Microservices, public validation |
| RS512     | RSA  | Slower | Microservices, higher security   |

**Use HMAC for single services**, RSA for microservices.

### AES Key Sizes

| Key Size | Security | Performance     |
| -------- | -------- | --------------- |
| AES-128  | Good     | Fastest         |
| AES-192  | Better   | Fast            |
| AES-256  | Best     | Still very fast |

**Use AES-256 unless you have size constraints**.

## Error Handling

All crypto functions return errors with context:

```go
hash, err := crypto.HashPasswordBcrypt("")
if err != nil {
    // Error: "password cannot be empty"
    log.Error("Hashing failed", logger.Err(err))
}

claims, err := crypto.ParseJWT(invalidToken)
if err != nil {
    if errors.Is(err, errors.CodeTokenExpired) {
        // Handle expired token
    } else if errors.Is(err, errors.CodeInvalidToken) {
        // Handle invalid token
    }
}
```

## Performance

- **Bcrypt**: ~100-200 hashes/second (intentionally slow)
- **Argon2**: ~50-100 hashes/second (configurable)
- **JWT**: 100,000+ tokens/second
- **AES-GCM**: Millions of operations/second

Password hashing is intentionally slow to prevent brute force attacks.

## Next Steps

- **[Authentication](../authentication/getting-started)** - Complete auth system
- **[Logger](./logger)** - Structured logging
- **[Security Best Practices](../security/best-practices)** - Security guide

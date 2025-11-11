# Authentication Service Example

A comprehensive example demonstrating the use of the go-infra authentication component with microservice-style implementation.

## 🎯 What's Demonstrated

This example showcases:

### ✅ Authentication Features

- ✅ **Sign Up** - User registration with identifier/credential pattern
- ✅ **Sign In** - Authentication with JWT tokens
- ✅ **Profile Management** - Get and update user profile
- ✅ **Password Change** - Secure password update flow
- ✅ **Protected Routes** - JWT middleware integration
- ✅ **Role-Based Access** - Admin endpoint example

### ✅ Technical Features

- ✅ **Swagger Documentation** - Complete API documentation with Swagger UI
- ✅ **SQL Migrations** - Database schema management with Goose
- ✅ **Dependency Injection** - Using Uber FX
- ✅ **PostgreSQL** - Production-ready database setup
- ✅ **JWT Authentication** - Secure token-based auth
- ✅ **Microservice Pattern** - Separation of User/Identifier/Credential/Profile

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- `swag` CLI tool (for swagger generation)

### 1. Install Dependencies

```bash
go mod download
```

### 2. Setup Environment

Copy the example environment file:

```bash
cp .env.example .env
```

Update `.env` with your configuration:

```bash
# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=auth_service_db
POSTGRES_SSLMODE=disable

# JWT
JWT_SECRET=your-super-secret-key-change-in-production
JWT_ISSUER=authentication-service
JWT_AUDIENCE=web-app
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d
```

### 3. Create Database

```bash
createdb auth_service_db
```

### 4. Run Migrations

```bash
go run cmd/migration/main.go
```

This will create:

- `ms_users` - Core user entity
- `ms_user_identifiers` - Username, email, phone identifiers
- `ms_user_credentials` - Password and credential storage
- `ms_user_profiles` - User profile information

### 5. Run the Service

```bash
go run cmd/app/main.go
```

The service will start on `http://localhost:8081`

### 6. Access Swagger UI

Open your browser:

```
http://localhost:8081/swagger/index.html
```

## 📚 Documentation

- **[Swagger Guide](docs/SWAGGER_GUIDE.md)** - Complete guide for API documentation
- **[Migration Guide](docs/MIGRATION_GUIDE.md)** - Database migration best practices
- **[Architecture](docs/FRAMEWORK_VS_IMPLEMENTATION.md)** - Component architecture
- **[Switching Auth](docs/SWITCHING_AUTH_APPROACHES.md)** - How to switch between auth patterns

## 🔌 API Endpoints

### Public Endpoints

#### Sign Up

```bash
POST /api/v1/auth/signup
Content-Type: application/json

{
  "identifier": {
    "scheme": "username",
    "value": "john_doe"
  },
  "credential": {
    "scheme": "basic",
    "value": "secure_password123"
  },
  "firstname": "John",
  "lastname": "Doe",
  "email": "john@example.com"
}
```

#### Sign In

```bash
POST /api/v1/auth/signin
Content-Type: application/json

{
  "identifier": {
    "scheme": "username",
    "value": "john_doe"
  },
  "credential": {
    "scheme": "basic",
    "value": "secure_password123"
  }
}
```

Response:

```json
{
	"userId": "123e4567-e89b-12d3-a456-426614174000",
	"username": "john_doe",
	"email": "john@example.com",
	"token": {
		"value": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		"scheme": "bearer",
		"expiresAt": "2024-01-01T01:00:00Z",
		"type": "access"
	}
}
```

### Protected Endpoints (Require Authentication)

Add JWT token to headers:

```
Authorization: Bearer <your-access-token>
```

#### Get Profile

```bash
GET /api/v1/auth/profile
Authorization: Bearer <token>
```

#### Update Profile

```bash
PUT /api/v1/auth/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "firstname": "John",
  "lastname": "Smith",
  "email": "john.smith@example.com"
}
```

#### Change Password

```bash
PUT /api/v1/auth/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "oldCredential": {
    "scheme": "basic",
    "value": "old_password"
  },
  "newCredential": {
    "scheme": "basic",
    "value": "new_password123"
  }
}
```

#### User Dashboard

```bash
GET /api/v1/protected/dashboard
Authorization: Bearer <token>
```

#### User Settings

```bash
GET /api/v1/protected/settings
Authorization: Bearer <token>

PUT /api/v1/protected/settings
Authorization: Bearer <token>
Content-Type: application/json

{
  "theme": "dark",
  "language": "en"
}
```

### Admin Endpoints (Require Admin Role)

#### Admin Dashboard

```bash
GET /api/v1/admin/dashboard
Authorization: Bearer <token>
```

> **Note:** This endpoint checks for "admin" role in JWT claims

## 🧪 Testing with cURL

### 1. Sign Up

```bash
curl -X POST http://localhost:8081/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": {"scheme": "username", "value": "testuser"},
    "credential": {"scheme": "basic", "value": "password123"},
    "firstname": "Test",
    "lastname": "User",
    "email": "test@example.com"
  }'
```

### 2. Sign In and Save Token

```bash
TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/auth/signin \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": {"scheme": "username", "value": "testuser"},
    "credential": {"scheme": "basic", "value": "password123"}
  }' | jq -r '.token.value')

echo "Token: $TOKEN"
```

### 3. Access Protected Endpoint

```bash
curl http://localhost:8081/api/v1/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Get Dashboard

```bash
curl http://localhost:8081/api/v1/protected/dashboard \
  -H "Authorization: Bearer $TOKEN"
```

## 🏗️ Project Structure

```
authentication-service/
├── cmd/
│   ├── app/            # Main application entry point
│   └── migration/      # Database migration runner
├── config/             # Configuration loading
├── db/
│   └── migrations/
│       └── goose-migrate/  # SQL migration files
├── docs/               # Generated swagger docs
├── internal/
│   ├── microservice-auth/  # Microservice auth implementation
│   └── shared/
│       ├── app/        # Application setup
│       ├── data/       # Database context
│       └── handlers/   # HTTP handlers with swagger annotations
├── .env.example        # Example environment variables
├── go.mod
└── README.md
```

## 🔧 Configuration

### Database Configuration

```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=auth_service_db
POSTGRES_SSLMODE=disable
POSTGRES_MAX_OPEN_CONNS=25
POSTGRES_MAX_IDLE_CONNS=5
POSTGRES_CONN_MAX_LIFETIME=5m
```

### JWT Configuration

```bash
JWT_SECRET=change-this-to-a-secure-secret-key
JWT_ISSUER=authentication-service
JWT_AUDIENCE=web-app
JWT_ACCESS_EXPIRY=15m      # Access token lifetime
JWT_REFRESH_EXPIRY=168h    # Refresh token lifetime (7 days)
```

### Server Configuration

```bash
HTTP_PORT=8081
HTTP_BASE_PATH=/api/v1
HTTP_READ_TIMEOUT=30s
HTTP_WRITE_TIMEOUT=30s
```

## 🗄️ Database Schema

### Tables

#### users

Core user entity:

- `id` - UUID primary key
- `user_type` - User role (user, admin, etc.)
- `status` - Account status (active, inactive, suspended, banned)
- `last_login_at` - Last login timestamp
- `created_at`, `updated_at`, `deleted_at` - Timestamps

#### user_identifiers

How users can be identified:

- `id` - UUID primary key
- `user_id` - Foreign key to users
- `scheme` - Identifier type (username, email, phone)
- `identifier` - The actual identifier value
- `verified` - Whether identifier is verified
- `verified_at` - Verification timestamp

#### user_credentials

Authentication credentials:

- `id` - UUID primary key
- `user_id` - Foreign key to users
- `scheme` - Credential type (basic, oauth, apikey)
- `credential` - Hashed credential (e.g., bcrypt password)
- `expires_at` - Expiration timestamp

#### user_profiles

User profile information:

- `id` - UUID primary key
- `user_id` - Foreign key to users (unique)
- `firstname`, `lastname` - Name fields
- `email` - Email address
- `birthday` - Date of birth
- `locale` - Language preference
- `metadata` - JSONB for custom fields

## 🔄 Migration Management

### Run Migrations

```bash
go run cmd/migration/main.go
```

### Create New Migration

```bash
cd db/migrations/goose-migrate
goose create add_new_feature sql
```

### Migration Status

```bash
goose -dir db/migrations/goose-migrate postgres "connection-string" status
```

### Rollback

```bash
goose -dir db/migrations/goose-migrate postgres "connection-string" down
```

For detailed migration guidance, see [Migration Guide](docs/MIGRATION_GUIDE.md)

## 📝 Swagger Documentation

### Generate/Update Swagger

```bash
swag init --parseDependency --parseInternal -g cmd/app/main.go -o docs
```

### Install swag

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

For detailed swagger guidance, see [Swagger Guide](docs/SWAGGER_GUIDE.md)

## 🎓 Learning Resources

### Understanding the Architecture

1. **Framework vs Implementation**

   - Read: [Framework vs Implementation](docs/FRAMEWORK_VS_IMPLEMENTATION.md)
   - The component provides infrastructure, you provide business logic

2. **Microservice Auth Pattern**

   - Separation of concerns: User / Identifier / Credential / Profile
   - Flexible identification schemes (username, email, phone)
   - Multiple credential schemes (password, OAuth, API keys)

3. **Identifier/Credential Pattern**
   - Decouple "who you are" from "how you prove it"
   - Support multiple identifiers per user
   - Support multiple credential types

### Key Concepts

**Identifier**: How a user identifies themselves

- Scheme: `username`, `email`, `phone`
- Value: The actual identifier

**Credential**: How a user proves their identity

- Scheme: `basic` (password), `oauth`, `apikey`
- Value: The credential data (hashed)

**Profile**: User information (name, email, birthday, etc.)

## 🔐 Security Best Practices

### 1. JWT Secrets

- Use strong, random secrets in production
- Rotate secrets periodically
- Store secrets in secure vault (not in code)

### 2. Password Security

- Minimum 8 characters
- Passwords are hashed with bcrypt (cost factor 12)
- Never log or expose passwords

### 3. Token Expiry

- Access tokens: Short-lived (15 minutes)
- Refresh tokens: Longer-lived (7 days)
- Implement token refresh flow in production

### 4. HTTPS

- Always use HTTPS in production
- Redirect HTTP to HTTPS

### 5. Rate Limiting

- Implement rate limiting for auth endpoints
- Prevent brute force attacks

## 🧩 Integration Examples

### Using in Your Application

```go
import (
    authComponent "github.com/phatnt199/go-infra/pkg/component/authentication"
    "github.com/phatnt199/go-infra/pkg/crypto"
)

// Create component with default implementation
authComp, err := authComponent.NewComponentWithDefaultImplementation(
    db,
    jwtConfig,
    authComponent.WithLogger(logger),
)

// Get services
authService := authComp.GetAuthService()
jwtMiddleware := authComp.GetJWTMiddleware()

// Register routes
router.POST("/auth/signin", handleSignIn(authService))
router.GET("/protected", jwtMiddleware.Handle(), handleProtected)
```

### Custom User Provider

If you have your own user model:

```go
type MyUserProvider struct {
    db *gorm.DB
}

func (p *MyUserProvider) GetUserByIdentifier(ctx context.Context, scheme, value string) (*contracts.UserAuthInfo, error) {
    // Your implementation
}

// Create component with custom provider
authComp := authComponent.NewComponent(
    &MyUserProvider{db: db},
    authComponent.WithLogger(logger),
)
```

## 🐛 Troubleshooting

### "Table already exists" error

- **Cause**: Both AutoMigrate and SQL migrations enabled
- **Fix**: See [Migration Guide](docs/MIGRATION_GUIDE.md)

### Swagger paths empty

- **Cause**: Handlers not scanned by swaggo
- **Fix**: See [Swagger Guide](docs/SWAGGER_GUIDE.md)

### JWT validation fails

- **Cause**: Token expired or invalid secret
- **Fix**: Check JWT_SECRET matches, verify token expiry

### Database connection fails

- **Cause**: PostgreSQL not running or wrong credentials
- **Fix**: Check PostgreSQL is running, verify .env settings

## 📊 Production Considerations

### 1. Database

- Use connection pooling
- Set appropriate timeout values
- Enable query logging for debugging
- Regular backups

### 2. Security

- Use environment variables for secrets
- Enable CORS appropriately
- Implement rate limiting
- Add request logging

### 3. Monitoring

- Add health check endpoints
- Monitor JWT token generation/validation
- Track authentication failures
- Set up alerts for suspicious activity

### 4. Performance

- Add caching for frequently accessed data
- Optimize database queries
- Use CDN for static assets
- Implement pagination

## 🤝 Contributing

This is an example project. For contributions to the authentication component itself, see the main go-infra repository.

## 📄 License

MIT License - See LICENSE file for details

## 🔗 Related

- [go-infra](../../) - Main framework
- [Authentication Component](../../pkg/component/authentication/) - Component source
- [Other Examples](../) - More examples

## 📞 Support

- Issues: [GitHub Issues](https://github.com/phatnt199/go-infra/issues)
- Discussions: [GitHub Discussions](https://github.com/phatnt199/go-infra/discussions)

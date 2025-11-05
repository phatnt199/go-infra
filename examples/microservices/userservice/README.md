# User Service - Complete Implementation Guide

This document provides a complete guide to the userservice implementation with all required features.

## ✅ Completed Features

### 1. Database Schema with Soft Delete

All tables implement soft delete using the `deleted_at` column. The schema includes:

- **users** table with full user management
- **user_identifiers** (1-1 relationship) for username/phone authentication
- **user_credentials** (1-1 relationship) for password management
- **user_profiles** (1-1 relationship) for user profile data

All datetime columns use PostgreSQL TIMESTAMP which automatically formats to ISO 8601 (RFC3339) when queried.

### 2. BaseEntity

Created in `internal/shared/data/models/base_entity.go` with:

- `id` (UUID, auto-generated)
- `createdAt` (timestamp, auto-set)
- `modifiedAt` (timestamp, auto-updated)
- `deletedAt` (timestamp for soft delete)

### 3. Service Layer

Complete implementation with:

- **Authentication**: SignUp, SignIn, ChangePassword
- **User Management**: GetUserFullDetails, UpdateProfile, DisableUser, EnableUser
- **Password Hashing**: Using bcrypt via `pkg/crypto.Hasher`
- **JWT Tokens**: Generated via `pkg/crypto.JWTManager`

### 4. HTTP API Endpoints

#### Auth Endpoints

- `POST /api/v1/auth/signup` - Register new user
- `POST /api/v1/auth/signin` - Authenticate and get JWT token
- `POST /api/v1/auth/change-password/:id` - Change user password

#### User Endpoints

- `GET /api/v1/users/:id` - Get full user details (no password)
- `PUT /api/v1/users/:id/profile` - Update user profile
- `POST /api/v1/users/:id/disable` - Deactivate user account
- `POST /api/v1/users/:id/enable` - Activate user account

### 5. Database Seeding

Seeder implemented with 5 default users:

| Username   | Password       | Role                 |
| ---------- | -------------- | -------------------- |
| superadmin | SuperAdmin123! | System Administrator |
| admin      | Admin123!      | Administrator        |
| maintainer | Maintainer123! | Maintainer           |
| manager    | Manager123!    | Manager              |
| customer   | Customer123!   | Customer             |

All users use the "username" identifier scheme (not phoneNumber).

## 📁 Project Structure

```
userservice/
├── cmd/
│   ├── app/main.go              # Main application
│   └── seed/main.go             # Database seeder
├── db/migrations/goose-migrate/
│   └── 00001_create_users_table.sql
├── internal/
│   ├── shared/
│   │   ├── data/models/
│   │   │   └── base_entity.go
│   │   └── utils/
│   │       └── password.go
│   └── users/
│       ├── contracts/           # Interfaces
│       ├── data/
│       │   ├── datamodels/     # GORM models
│       │   └── repositories/   # Data access layer
│       ├── dtos/v1/            # Request/Response DTOs
│       ├── features/           # HTTP handlers
│       │   ├── auth/
│       │   └── profile/
│       ├── mappers/            # DTO converters
│       ├── models/             # Domain models
│       ├── seeder/             # Database seeding
│       ├── services/           # Business logic
│       └── users_fx.go         # Dependency injection
└── IMPLEMENTATION_SUMMARY.md   # Detailed summary
```

## 🗄️ Database Design

### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    status INT NOT NULL DEFAULT 100,  -- 100=ACTIVATED, 101=DEACTIVATED
    user_type VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    activated_at TIMESTAMP,
    last_login_at TIMESTAMP,
    parent_id UUID REFERENCES users(id),
    valid_from TIMESTAMP,
    valid_to TIMESTAMP
);
```

### User Identifiers (1-1 with users)

```sql
CREATE TABLE user_identifiers (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    modified_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    scheme VARCHAR(50) DEFAULT 'username',  -- 'username' or 'phoneNumber'
    identifier VARCHAR(255) NOT NULL,       -- actual username/phone
    verified BOOLEAN DEFAULT TRUE,
    details JSONB DEFAULT '{}'
);
```

### User Credentials (1-1 with users)

```sql
CREATE TABLE user_credentials (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    modified_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    scheme VARCHAR(50) DEFAULT 'basic',  -- 'basic' or 'master-basic'
    credential VARCHAR(255) NOT NULL,    -- hashed password
    details JSONB DEFAULT '{}'
);
```

### User Profiles (1-1 with users)

```sql
CREATE TABLE user_profiles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    modified_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    firstname VARCHAR(255),
    lastname VARCHAR(255),
    birthday DATE,
    locale VARCHAR(10) DEFAULT 'en_US',
    details JSONB DEFAULT '{}'
);
```

## 🚀 How to Run

### 1. Set up Database

Ensure PostgreSQL is running and create the database:

```bash
createdb userservice_db
```

### 2. Configure Environment

Update your `.env` or configuration file with database credentials.

### 3. Run Migrations

```bash
cd db/migrations/goose-migrate
goose postgres "postgres://user:password@localhost:5432/userservice_db?sslmode=disable" up
```

### 4. Seed Database

```bash
go run cmd/seed/main.go
```

This will create 5 default users. You'll see output like:

```
Successfully seeded user: superadmin (password: SuperAdmin123!)
Successfully seeded user: admin (password: Admin123!)
...
```

### 5. Start Application

```bash
go run cmd/app/main.go
```

## 📝 API Examples

### Sign Up

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "SecurePass123!",
    "firstname": "John",
    "lastname": "Doe",
    "locale": "en_US"
  }'
```

### Sign In

```bash
curl -X POST http://localhost:8080/api/v1/auth/signin \
  -H "Content-Type: application/json" \
  -d '{
    "username": "superadmin",
    "password": "SuperAdmin123!"
  }'
```

### Get User Details

```bash
curl -X GET http://localhost:8080/api/v1/users/{user-id} \
  -H "Authorization: Bearer {jwt-token}"
```

### Update Profile

```bash
curl -X PUT http://localhost:8080/api/v1/users/{user-id}/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {jwt-token}" \
  -d '{
    "firstname": "Jane",
    "lastname": "Smith",
    "locale": "vi_VN"
  }'
```

### Change Password

```bash
curl -X POST http://localhost:8080/api/v1/auth/change-password/{user-id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {jwt-token}" \
  -d '{
    "oldPassword": "SuperAdmin123!",
    "newPassword": "NewSecurePass456!"
  }'
```

### Disable User

```bash
curl -X POST http://localhost:8080/api/v1/users/{user-id}/disable \
  -H "Authorization: Bearer {jwt-token}"
```

### Enable User

```bash
curl -X POST http://localhost:8080/api/v1/users/{user-id}/enable \
  -H "Authorization: Bearer {jwt-token}"
```

## 🔐 Security Features

- ✅ Passwords hashed with bcrypt (cost: default)
- ✅ JWT token-based authentication
- ✅ Soft delete for data recovery
- ✅ SQL injection protection via GORM
- ✅ Input validation via go-playground/validator

## 📊 Response Formats

### Success Response (Auth)

```json
{
	"userId": "uuid",
	"accessToken": "jwt-token",
	"refreshToken": "jwt-refresh-token"
}
```

### Success Response (User Details)

```json
{
	"user": {
		"id": "uuid",
		"status": 100,
		"userType": "SYSTEM",
		"activatedAt": "2024-01-01T00:00:00Z",
		"createdAt": "2024-01-01T00:00:00Z",
		"modifiedAt": "2024-01-01T00:00:00Z"
	},
	"identifier": {
		"id": "uuid",
		"scheme": "username",
		"identifier": "superadmin",
		"verified": true
	},
	"profile": {
		"id": "uuid",
		"firstname": "Super",
		"lastname": "Administrator",
		"locale": "en_US",
		"createdAt": "2024-01-01T00:00:00Z"
	}
}
```

### Error Response

```json
{
	"success": false,
	"message": "Error description"
}
```

## 🎯 Key Design Decisions

1. **Separate Tables for Identity**: Following your Loopback 4 design pattern, user authentication (identifier), credentials, and profile are in separate tables with explicit `user_id` foreign keys.

2. **ISO DateTime Format**: All datetime fields use PostgreSQL TIMESTAMP and are formatted to ISO 8601 (RFC3339) in API responses.

3. **Soft Delete**: Implemented via GORM's `DeletedAt` field - deleted records are marked, not removed.

4. **No Roles/Permissions**: As requested, role-based access control is not implemented yet.

5. **Scheme Support**: Both identifier and credential tables support different schemes for future extensibility (e.g., phoneNumber, master-basic).

## ⚠️ Important Notes

- **Default Passwords**: The seeded users have simple passwords for testing. Change them in production!
- **JWT Secret**: Ensure you configure a strong JWT secret in your configuration.
- **Database Migrations**: Always backup before running migrations in production.
- **Soft Delete**: Remember that deleted records still exist in the database. Implement a cleanup job if needed.

## 📚 Next Steps

To add more features:

1. Implement roles and permissions
2. Add email verification
3. Implement password reset flow
4. Add 2FA support
5. Create admin panel for user management
6. Add audit logging
7. Implement rate limiting
8. Add user search and filtering

## 🐛 Troubleshooting

### Migration Fails

- Check database connection string
- Ensure PostgreSQL is running
- Verify database exists

### Seeder Error "user already exists"

- The seeder skips existing users
- Drop and recreate the database for a fresh start

### JWT Token Invalid

- Check JWT configuration
- Ensure token hasn't expired
- Verify token is sent in Authorization header

## 📞 Support

For issues or questions, please refer to:

- `IMPLEMENTATION_SUMMARY.md` for detailed technical information
- Service code comments for inline documentation
- Migration files for database schema details

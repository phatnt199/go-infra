---
sidebar_position: 1
slug: /
---

# Introduction

**go-infra** is a comprehensive, production-ready infrastructure framework for building Go applications. It provides everything you need to build scalable, maintainable backend services without the boilerplate.

## Why go-infra?

Building production applications requires more than just writing business logic. You need:

- **Database integration** with migrations and repositories
- **Authentication & authorization** with JWT and role-based access
- **Structured logging** with context propagation
- **HTTP server** with routing and middleware
- **Configuration management** with environment support
- **Testing utilities** for reliable tests
- **Security tools** for encryption and hashing

go-infra provides all of this out of the box, following Go's idiomatic patterns.

## Key Features

### 🚀 Rapid Development

- **Modular architecture** - Build REST APIs efficiently with reusable components
- **Module system** - Plug and play components
- **Dependency injection** - Powered by Uber Fx
- **Auto-configuration** - Smart defaults that just work

### 🔒 Security First

- **JWT authentication** with HMAC and RSA support
- **Password hashing** with Bcrypt and Argon2
- **AES encryption** for sensitive data
- **Role-based access control**

### 📊 Production Ready

- **Structured logging** with Zap
- **OpenTelemetry** metrics and tracing
- **Health checks** for monitoring
- **Database migrations** with Goose and Go-Migrate

### 🏗️ Clean Architecture

- **CQRS pattern** - Command/Query separation with mediator pattern
- **Event sourcing** - Event-driven architecture support (experimental)
- **Repository pattern** - Clean data access layer
- **Domain-driven design** helpers

## Quick Example

Here's how easy it is to build a complete REST API:

```go
package main

import (
    "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
    "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
)

func main() {
    app := fxapp.NewApplicationBuilder().
        ProvideModule(fiber_adapter.Module).
        ProvideModule(gorm.Module).
        Build()

    app.Run()
}
```

That's it! You now have:

- ✅ HTTP server with Fiber
- ✅ PostgreSQL database with GORM
- ✅ Structured logging
- ✅ Configuration management
- ✅ Graceful shutdown

## What's Inside?

### Core Packages

| Package                      | Purpose                                         |
| ---------------------------- | ----------------------------------------------- |
| `adapter/fxapp`              | Application bootstrap with dependency injection |
| `adapter/http/fiber_adapter` | HTTP server with Fiber                          |
| `infra/postgres/gorm`        | PostgreSQL integration with GORM                |
| `component/authentication`   | Complete auth system                            |
| `logger`                     | Structured logging                              |
| `crypto`                     | Security utilities                              |
| `core/cqrs`                  | CQRS pattern with commands and queries          |
| `otel`                       | OpenTelemetry metrics (experimental)            |
| `mapper`                     | Object-to-object mapping with generics          |
| `reflection`                 | Reflection utilities for advanced scenarios     |
| `validator`                  | Struct validation wrapper                       |
| `utils`                      | Pagination, errors, and common utilities        |
| `migration`                  | Database migrations with Goose                  |

## Who Should Use This?

go-infra is perfect for:

- **Backend engineers** building REST APIs
- **Teams** wanting consistent architecture
- **Startups** needing rapid development
- **Microservices** requiring standardization
- **Anyone** who wants less boilerplate

## Next Steps

- **[Quick Start](./getting-started/quick-start)** - Build your first API in 5 minutes
- **[Architecture](./core-concepts/architecture)** - Understand how it works
- **[Examples](./examples/users-api)** - Learn from real code

## Coming Soon

The following features are planned for future releases:

- 🔄 **Redis integration** - Caching and session storage
- 📮 **Message queues** - RabbitMQ, Kafka support
- 🔍 **Advanced observability** - Distributed tracing with OpenTelemetry
- 🗄️ **More databases** - MySQL, MongoDB support

## Philosophy

go-infra follows these principles:

1. **Convention over configuration** - Smart defaults, customize when needed
2. **Explicit over magic** - Clear, understandable code
3. **Composable** - Use what you need, ignore the rest
4. **Production-first** - Built for real-world applications
5. **Developer experience** - Easy to use, hard to misuse

Ready to get started? Let's build something awesome! 🚀

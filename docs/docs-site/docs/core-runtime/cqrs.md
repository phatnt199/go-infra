---
title: CQRS Patterns
description: Command and query segregation support in go-infra.
---

Highlights command dispatching, query models, and consistency boundaries.

## Sections

- Command contracts
- Handler registration (direct vs Fx)
- Query side patterns
- Transaction & unit-of-work notes

> Placeholder: Implementation examples to follow.

## Command Example

```go
type CreateUserCommand struct { Email string }
type CreateUserHandler struct { repo UserRepo }
func (h CreateUserHandler) Handle(ctx context.Context, c CreateUserCommand) error { return h.repo.Create(ctx, c.Email) }
```

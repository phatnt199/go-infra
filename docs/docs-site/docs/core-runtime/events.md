---
title: Domain & Integration Events
description: Publishing, subscribing, and handling events in go-infra.
---

Outlines event abstractions and how they integrate with messaging and projections.

## Sections

- Event interfaces & envelopes
- Event bus vs domain events
- Subscription & handler registration
- Testing strategies

> Placeholder content. Code snippets will be generalized from example services.

## Publishing & Handling

```go
type EventPublisher interface { Publish(ctx context.Context, evt any) error }
type Handler interface { Handle(ctx context.Context, evt any) error }
```

Register handler:

```go
bus.Add(UserRegistered{}, handler)
```

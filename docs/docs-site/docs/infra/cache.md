---
title: Cache Layer
description: Abstractions for caching and pluggable backends.
---

## Sections

- Interface & contracts
- In-memory vs Redis
- Fx wiring
- Cache invalidation patterns

> Placeholder.

## Interface

```go
type Cache interface { Get(ctx context.Context, key string, dst any)(bool,error); Set(ctx context.Context, key string, v any, ttl time.Duration) error }
```

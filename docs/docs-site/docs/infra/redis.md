---
title: Redis Adapter
description: Caching and data structures via Redis.
---

## Sections

- Config & options
- Direct client usage
- Fx provision
- TTL & serialization considerations

> Placeholder.

## Client Usage

```go
client := redis.NewClient(&redis.Options{Addr:"localhost:6379"})
client.Set(ctx, "k", "v", time.Hour)
```

---
title: Goose Migrations
description: Using goose with go-infra.
---

## Sections

- Configuration & directory layout
- Running migrations (CLI & programmatic)
- Fx integration points
- Rollback & safety practices

> Placeholder.

## Runner

```go
runner := goose.NewGoosePostgres(cfg, db, log)
_ = runner.Up(ctx, 0)
```

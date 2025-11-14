---
title: GoMigrate Integration
description: Managing migrations with gomigrate implementation.
---

## Sections

- Setup & config structs
- Applying migrations
- Comparing to goose
- Fx usage pattern

> Placeholder.

## Runner

```go
runner, _ := gomigrate.NewGoMigratorPostgres(cfg, db, log)
runner.Up(ctx, 0)
```

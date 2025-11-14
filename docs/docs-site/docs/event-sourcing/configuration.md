---
title: Event Store Configuration
description: Configuring event sourcing components.
---

## Sections

- Config structs & environment
- Tuning performance
- Fx vs direct wiring

> Placeholder.

## Config Sample

```go
type Config struct { SnapshotFrequency int64 `json:"snapshotFrequency" validate:"required,gte=0"` }
```

---
title: Storage Adapter
description: File/object storage integration.
---

## Sections

- Abstraction interfaces
- Local vs cloud implementations
- Fx wiring
- Streaming uploads/downloads

> Placeholder.

## ObjectStore Interface

```go
type ObjectStore interface { Put(ctx context.Context, path string, r io.Reader, meta map[string]string) error }
```

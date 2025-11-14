---
title: Queue Adapter
description: Working with asynchronous queues.
---

## Sections

- Enqueue & dequeue APIs
- Reliability patterns
- Fx module usage

> Placeholder.

## Queue API

```go
type Queue interface { Enqueue(ctx context.Context, msg any) error }
```

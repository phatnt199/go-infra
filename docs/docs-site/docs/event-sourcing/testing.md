---
title: Testing Event Sourced Components
description: Patterns for verifying event streams, projections, and handlers.
---

## Sections

- In-memory test doubles
- Stream assertions
- Projection consistency tests

> Placeholder.

## Assertion

```go
events := store.Events(aggID)
require.Equal(t, "UserCreated", events[0].Name)
```

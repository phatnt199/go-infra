---
title: Projections & Read Models
description: Building and updating projections.
---

## Sections

- Projection publisher
- Checkpoint repository
- Rebuild strategies

> Placeholder.

## Publisher Snippet

```go
for _, pj := range p.projections { _ = pj.ProcessEvent(ctx, evt) }
```

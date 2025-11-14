---
title: Customizing Components
description: Overriding defaults and extending behaviors.
---

## Sections

- Replacing implementations via interfaces
- Configuration layering
- Fx override patterns

> Placeholder.

## Override Example

```go
fx.New(
	CoreModule,
	fx.Provide(NewDefaultCache),
	fx.Replace(func() Cache { return NewRedisCache(customOpts) }),
)
```

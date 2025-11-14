---
title: Serialization Strategies
description: Serializer abstractions, formats, and customization.
---

## Sections

- Default serializers
- Custom codec injection
- Handling versioning

> Placeholder.

## Custom Serializer

```go
type ProtoSerializer struct{}
func (ProtoSerializer) Marshal(v any) ([]byte, error) { return []byte{}, nil }
func (ProtoSerializer) Unmarshal(b []byte, v any) error { return nil }
```

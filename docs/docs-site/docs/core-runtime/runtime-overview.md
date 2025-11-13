---
title: Runtime Overview
description: Overview of the go-infra core runtime primitives and lifecycle.
---

This page introduces the runtime layer of `go-infra`: how modules, lifecycle hooks, and Fx integration compose your application.

## Goals

- Standardize bootstrapping
- Provide composable modules
- Support direct usage without Fx when desired

## Contents

1. Module registration pattern
2. Lifecycle and shutdown hooks
3. Configuration injection flow
4. Interplay between `pkg/core` and adapters

> Placeholder: Detailed examples adapted from examples will appear here with renamed package paths to keep documentation neutral.

## Module Registration

```go
import (
	"go.uber.org/fx"
	corejson "github.com/phatnt199/go-infra/pkg/core/serializer/json"
)

var CoreModule = fx.Module(
	"corefx",
	fx.Provide(
		corejson.NewDefaultJsonSerializer,
		corejson.NewDefaultEventJsonSerializer,
		corejson.NewDefaultMessageJsonSerializer,
		corejson.NewDefaultMetadataJsonSerializer,
	),
)
```

## Lifecycle Hooks

```go
type Runner struct{}
func NewRunner(lc fx.Lifecycle) *Runner {
	r := &Runner{}
	lc.Append(fx.Hook{OnStart: func(ctx context.Context) error { return nil }, OnStop: func(ctx context.Context) error { return nil }})
	return r
}
```

## Direct Usage

```go
s := corejson.NewDefaultJsonSerializer()
evtSer := corejson.NewDefaultEventJsonSerializer()
```

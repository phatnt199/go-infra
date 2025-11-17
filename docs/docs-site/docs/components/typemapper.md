---
title: Typemapper
description: Runtime registry and helpers for discovering types at runtime.
---

# Typemapper

`pkg/reflection/typemapper` discovers and registers runtime types so other components can look up types by name, create instances, or validate interface implementations without hard-coded type references.

## Automatic discovery

- The package runs `discoverTypes` during `init()` by invoking the linked `reflect.typelinks` and `reflect.resolveTypeOff` helpers (`unsafe_types.go`).
- It walks every symbol, registers pointer-to-struct types, and stores both the pointer and dereferenced struct under multiple keys:
  - `types map[string][]reflect.Type` - keyed by both full `pkg.Path.TypeName` and short names (`*TypeName`/`TypeName`).
  - `packages map[string]map[string][]reflect.Type` - indexed by package path then short name.
- Discovery enables finding types without manually wiring them, but the package also supports explicit registration when determinism or cross-version guarantees are required.

## Registry helpers

| Function                                                  | Description                                           |
| --------------------------------------------------------- | ----------------------------------------------------- |
| `RegisterType(typ reflect.Type)`                          | Register a type (short and full names) manually.      |
| `RegisterTypeWithKey(key string, typ reflect.Type)`       | Bind a type to a custom lookup key.                   |
| `GetAllRegisteredTypes() map[string][]reflect.Type`       | Returns live registry map (read-only reference).      |
| `TypeByName(string) reflect.Type`                         | Returns the first type registered under the name/key. |
| `TypesByName(string) []reflect.Type`                      | Returns every registered type for the key.            |
| `TypeByPackageName(pkgPath, name string) reflect.Type`    | Lookup by package path and short name.                |
| `TypesByPackageName(pkgPath, name string) []reflect.Type` | Same but returns all matches.                         |

## Naming utilities

- `GetFullTypeName`, `GetFullTypeNameByType`, `GetGenericFullTypeNameByT` — return `pkg.Path.TypeName` or `*pkg.Path.TypeName` strings.
- `GetTypeName`, `GetNonePointerTypeName`, `GetTypeNameByType`, `GetGenericTypeNameByT`, `GetGenericNonePointerTypeNameByT` — return short names (optionally without the pointer `*`).
- `GetSnakeTypeName`, `GetKebabTypeName` — convert pointer type names into `snake_case` or `kebab-case`.
- `GetPackageName(value interface{})` extracts the final path segment from the package path.

## Instance creation

| Function                                                                       | Description                                                                    |
| ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| `GenericInstanceByT[T any]() T`                                                | Returns an empty instance of the generic type `T`.                             |
| `InstanceByType(typ reflect.Type) interface{}`                                 | Creates an empty pointer (if `typ` is pointer) or zero value of `typ`.         |
| `InstanceByTypeName(name string)`                                              | Lookup a type by name in the registry and return a blank value.                |
| `InstancePointerByTypeName(name string)`                                       | Always returns a pointer, even when `name` resolved to a struct type.          |
| `InstanceByPackageName(pkgPath, name string)`                                  | Creates an instance for the type named `name` in package `pkgPath`.            |
| `EmptyInstanceByTypeNameAndImplementedInterface[TInterface any](name string)`  | Finds the first registered type with name `name` that implements `TInterface`. |
| `EmptyInstanceByTypeAndImplementedInterface[TInterface any](typ reflect.Type)` | Helper that resolves by short type name and interface.                         |

## Interface helpers

- `TypesImplementedInterface[TInterface any]()` returns every registered type implementing the interface.
- `TypesImplementedInterfaceWithFilterTypes[TInterface any](types []reflect.Type)` filters an existing slice.
- `TypeByNameAndImplementedInterface[TInterface any](typeName string)` finds the first type matching both name and interface.
- `GetGenericImplementInterfaceTypesT[T any]()` groups implementing types by registry key.
- `ImplementedInterfaceT[T any](obj interface{}) bool` checks whether `obj` implements `T`.

## Usage examples

```go
import "github.com/phatnt199/go-infra/pkg/reflection/typemapper"

type Service interface { Serve() }

typ := typemapper.TypeByNameAndImplementedInterface[Service]("*DefaultService")
if typ == nil {
    log.Fatal("service type not registered")
}

svc := typemapper.InstanceByType(typ).(Service)
svc.Serve()
```

```go
customType := reflect.TypeOf(MyCustom{})
typemapper.RegisterType(customType)
inst := typemapper.InstanceByTypeName("MyCustom")
```

## Warnings

- Discovery relies on `unsafe` and private runtime symbols (`reflect.typelinks`). The behavior may vary across Go versions or when the binary is stripped; use explicit registration where stability is vital.
- `types` and `packages` maps store slices of types — `TypeByName` simply returns the first entry, so use the plural helpers when multiple matches exist.
- The registry only seeds pointer-to-struct types automatically; other kinds require manual registration.

## See also

- `pkg/reflection/typemapper/type_mapper.go`
- `pkg/reflection/typemapper/unsafe_types.go`
- `components/mapper.md`

<!-- EOF -->

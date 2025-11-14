---
title: Reflection Utilities
description: Helpers for runtime type inspection.
---

# Reflection Utilities

`reflection/reflection_helper` offers a small set of safe helpers for examining and mutating struct fields and invoking methods at runtime. Unlike `typemapper`, it works directly on pointer or value instances without relying on runtime symbol discovery.

## Key capabilities

- `GetAllFields(typ reflect.Type)` — returns every `StructField` defined on `typ`, stepping through embedded fields when necessary. Accepts pointer types by automatically dereferencing them.
- Field getters/setters by index or name:
  - `GetFieldValueByIndex[T any](object T, index int) interface{}`
  - `GetFieldValueByName[T any](object T, name string) interface{}`
  - `SetFieldValueByIndex[T any](object T, index int, value interface{})`
  - `SetFieldValueByName[T any](object T, name string, value interface{})`
    Those helpers internally use `unsafe.Pointer` to read/write unexported fields. When you pass a pointer, the helper dereferences it before operating. When the field is read-only from `reflect`, the helper bypasses it via `reflect.NewAt` and `unsafe.Pointer`.
- Low-level helpers:
  - `GetFieldValue(reflect.Value) reflect.Value` / `SetFieldValue(reflect.Value, interface{})` — work with `reflect.Value` directly for cases where you already have a value handle.
  - `GetFieldValueFromMethodAndObject[T interface{}](object T, name string)` / `GetFieldValueFromMethodAndReflectValue(reflect.Value, name string)` — call zero-arg methods by name and return the result, accounting for pointer receivers if needed.
  - `SetValue[T any](data T, value interface{})` — rounds up a `reflect.Value` to write either pointer or value types.

## Path utilities

- `ObjectTypePath(obj any) string` / `TypePath[T any]() string` — return `package.Path.TypeName` for values or generic type parameters.
- `MethodPath(f interface{}) string` — converts a function pointer into the cleaned-up string `package:Receiver.Method` (strips the compiler `-fm` suffix and pointer markers).

## Usage guidance

1. Reflectively inspect a struct:

   ```go
   fields := reflectionHelper.GetAllFields(reflect.TypeOf(&MyStruct{}))
   for _, field := range fields {
   		 fmt.Println(field.Name, field.Type)
   }
   ```

2. Read or write an unexported field:

   ```go
   val := reflectionHelper.GetFieldValueByName(myObj, "secret")
   reflectionHelper.SetFieldValueByName(&myObj, "secret", newValue)
   ```

3. Call a method and get its first return value without needing concrete interfaces:

   ```go
   result := reflectionHelper.GetFieldValueFromMethodAndObject(myObj, "Compute")
   ```

## Safety notes

- The helpers intentionally use `unsafe.Pointer` when reading or writing unexported fields, so they should only be used in trusted internal contexts (e.g., testing or framework plumbing).
- Because they bypass Go's export rules, you should avoid exposing these helpers directly in public APIs unless absolutely necessary.

<!-- EOF -->

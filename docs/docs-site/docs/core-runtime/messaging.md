---
title: Messaging Layer
description: Using the messaging primitives for asynchronous workflows.
---

Explains abstractions for publishing, consuming, and integration with queues.

## Sections

- Message contracts
- Publisher/consumer APIs
- Fx module provisioning
- Error handling & retries

> Placeholder.

## Message Example

```go
type UserSyncMessage struct { UserID string; Version int }
queue.Enqueue(ctx, UserSyncMessage{UserID: id, Version:1})
```

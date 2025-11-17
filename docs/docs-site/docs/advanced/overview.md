---
sidebar_position: 1
---

# Advanced Features

Explore advanced go-infra features including CQRS, Event Sourcing, WebSockets, and more.

## Command Query Responsibility Segregation (CQRS)

CQRS separates read and write operations for better scalability and performance.

### Basic CQRS Implementation

```go
// Commands - Write operations
type CreateUserCommand struct {
    Name  string
    Email string
}

type UpdateUserCommand struct {
    ID    string
    Name  string
    Email string
}

// Queries - Read operations
type GetUserQuery struct {
    ID string
}

type ListUsersQuery struct {
    Page  int
    Limit int
}
```

### Command Handler

```go
// internal/cqrs/command_handler.go
package cqrs

import (
    "github.com/phatnt199/go-infra/pkg/logger"
    "myapp/internal/domain"
    "myapp/internal/repository"
)

type UserCommandHandler struct {
    repo   repository.UserRepository
    logger logger.Logger
}

func (h *UserCommandHandler) HandleCreateUser(cmd CreateUserCommand) (*domain.User, error) {
    h.logger.Infow("Handling CreateUser command", logger.Fields{
        "email": cmd.Email,
    })

    user := &domain.User{
        Name:  cmd.Name,
        Email: cmd.Email,
    }

    if err := h.repo.Create(user); err != nil {
        return nil, err
    }

    // Publish event
    h.eventBus.Publish("user.created", UserCreatedEvent{
        UserID: user.ID,
        Email:  user.Email,
    })

    return user, nil
}
```

### Query Handler

```go
// internal/cqrs/query_handler.go
package cqrs

type UserQueryHandler struct {
    readRepo repository.UserReadRepository
}

func (h *UserQueryHandler) HandleGetUser(query GetUserQuery) (*domain.User, error) {
    // Read from optimized read model
    return h.readRepo.FindByID(query.ID)
}

func (h *UserQueryHandler) HandleListUsers(query ListUsersQuery) ([]*domain.User, error) {
    return h.readRepo.List(query.Page, query.Limit)
}
```

## Event Sourcing

Store state changes as a sequence of events.

### Event Store

```go
// internal/eventsourcing/event.go
package eventsourcing

import "time"

type Event struct {
    ID            string
    AggregateID   string
    AggregateType string
    EventType     string
    EventData     interface{}
    Version       int
    Timestamp     time.Time
}

type EventStore interface {
    SaveEvent(event Event) error
    GetEvents(aggregateID string) ([]Event, error)
    GetEventsByType(eventType string) ([]Event, error)
}
```

### Aggregate Root

```go
// internal/domain/user_aggregate.go
package domain

type UserAggregate struct {
    ID      string
    Version int
    Events  []Event

    // Current state
    Name  string
    Email string
}

func (a *UserAggregate) ApplyEvent(event Event) {
    switch event.EventType {
    case "UserCreated":
        data := event.EventData.(UserCreatedEvent)
        a.Name = data.Name
        a.Email = data.Email
    case "UserUpdated":
        data := event.EventData.(UserUpdatedEvent)
        a.Name = data.Name
        a.Email = data.Email
    }
    a.Version++
}

func (a *UserAggregate) CreateUser(name, email string) {
    event := Event{
        EventType: "UserCreated",
        EventData: UserCreatedEvent{Name: name, Email: email},
    }
    a.Events = append(a.Events, event)
    a.ApplyEvent(event)
}
```

## WebSocket Support

Real-time bidirectional communication.

### WebSocket Server

```go
// internal/websocket/server.go
package websocket

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
    "github.com/phatnt199/go-infra/pkg/logger"
)

type WebSocketServer struct {
    clients map[*websocket.Conn]bool
    logger  logger.Logger
}

func NewWebSocketServer(logger logger.Logger) *WebSocketServer {
    return &WebSocketServer{
        clients: make(map[*websocket.Conn]bool),
        logger:  logger,
    }
}

func (s *WebSocketServer) HandleConnection(c *websocket.Conn) {
    defer func() {
        delete(s.clients, c)
        c.Close()
    }()

    s.clients[c] = true
    s.logger.Infow("New WebSocket connection", logger.Fields{
        "total": len(s.clients),
    })

    for {
        messageType, message, err := c.ReadMessage()
        if err != nil {
            s.logger.Errorw("Read error", logger.Fields{
                "error": err.Error(),
            })
            break
        }

        s.logger.Infow("Received message", logger.Fields{
            "message": string(message),
        })

        // Broadcast to all clients
        s.Broadcast(messageType, message)
    }
}

func (s *WebSocketServer) Broadcast(messageType int, message []byte) {
    for client := range s.clients {
        err := client.WriteMessage(messageType, message)
        if err != nil {
            s.logger.Errorw("Write error", logger.Fields{
                "error": err.Error(),
            })
            client.Close()
            delete(s.clients, client)
        }
    }
}
```

### Register WebSocket Route

```go
// main.go
func main() {
    app := fiber.New()
    appLog := defaultlogger.GetLogger()
    wsServer := websocket.NewWebSocketServer(appLog)

    app.Get("/ws", websocket.New(func(c *websocket.Conn) {
        wsServer.HandleConnection(c)
    }))

    app.Listen(":3000")
}
```

### Client Example

```html
<!DOCTYPE html>
<html>
	<body>
		<script>
			const ws = new WebSocket("ws://localhost:3000/ws");

			ws.onopen = () => {
				console.log("Connected");
				ws.send("Hello Server!");
			};

			ws.onmessage = (event) => {
				console.log("Received:", event.data);
			};
		</script>
	</body>
</html>
```

## gRPC Support

High-performance RPC framework.

### Protocol Buffer Definition

```protobuf
// proto/user.proto
syntax = "proto3";

package user;

option go_package = "github.com/yourorg/myapp/proto/user";

service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
}

message GetUserRequest {
  string id = 1;
}

message GetUserResponse {
  string id = 1;
  string name = 2;
  string email = 3;
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
}

message CreateUserResponse {
  string id = 1;
  string name = 2;
  string email = 3;
}
```

### Generate Code

```bash
# Install protoc compiler
brew install protobuf  # macOS
# or
apt-get install protobuf-compiler  # Linux

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate code
protoc --go_out=. --go-grpc_out=. proto/user.proto
```

### gRPC Server

```go
// internal/grpc/server.go
package grpc

import (
    "context"
    pb "myapp/proto/user"
    "myapp/internal/service"
)

type UserGRPCServer struct {
    pb.UnimplementedUserServiceServer
    userService *service.UserService
}

func NewUserGRPCServer(userService *service.UserService) *UserGRPCServer {
    return &UserGRPCServer{userService: userService}
}

func (s *UserGRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    user, err := s.userService.GetByID(req.Id)
    if err != nil {
        return nil, err
    }

    return &pb.GetUserResponse{
        Id:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}

func (s *UserGRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    user, err := s.userService.Create(req.Name, req.Email)
    if err != nil {
        return nil, err
    }

    return &pb.CreateUserResponse{
        Id:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}
```

### Start gRPC Server

```go
// main.go
import (
    "google.golang.org/grpc"
    pb "myapp/proto/user"
)

func main() {
    // Create gRPC server
    grpcServer := grpc.NewServer()

    // Register service
    userGRPCServer := grpc.NewUserGRPCServer(userService)
    pb.RegisterUserServiceServer(grpcServer, userGRPCServer)

    // Listen
    lis, _ := net.Listen("tcp", ":50051")
    log.Println("gRPC server listening on :50051")
    grpcServer.Serve(lis)
}
```

## Message Queue Integration

Async communication with RabbitMQ.

### RabbitMQ Publisher

```go
// internal/messaging/publisher.go
package messaging

import (
    "encoding/json"
    "github.com/streadway/amqp"
)

type Publisher struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func NewPublisher(url string) (*Publisher, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    channel, err := conn.Channel()
    if err != nil {
        return nil, err
    }

    return &Publisher{conn: conn, channel: channel}, nil
}

func (p *Publisher) Publish(exchange, routingKey string, message interface{}) error {
    body, err := json.Marshal(message)
    if err != nil {
        return err
    }

    return p.channel.Publish(
        exchange,
        routingKey,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}
```

### RabbitMQ Consumer

```go
// internal/messaging/consumer.go
package messaging

type Consumer struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func NewConsumer(url string) (*Consumer, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    channel, err := conn.Channel()
    if err != nil {
        return nil, err
    }

    return &Consumer{conn: conn, channel: channel}, nil
}

func (c *Consumer) Consume(queueName string, handler func([]byte) error) error {
    msgs, err := c.channel.Consume(
        queueName,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    go func() {
        for msg := range msgs {
            if err := handler(msg.Body); err != nil {
                log.Printf("Handler error: %v", err)
            }
        }
    }()

    return nil
}
```

## Redis Caching

Improve performance with caching.

```go
// internal/cache/redis.go
package cache

import (
    "context"
    "encoding/json"
    "time"
    "github.com/redis/go-redis/v9"
)

type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })
    return &RedisCache{client: client}
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := c.client.Get(ctx, key).Bytes()
    if err != nil {
        return err
    }
    return json.Unmarshal(data, dest)
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
    return c.client.Del(ctx, key).Err()
}
```

## Background Jobs

Process tasks asynchronously.

```go
// internal/jobs/worker.go
package jobs

import (
    "context"
    "time"
)

type Job interface {
    Execute(ctx context.Context) error
}

type Worker struct {
    jobQueue chan Job
    quit     chan bool
}

func NewWorker(queueSize int) *Worker {
    return &Worker{
        jobQueue: make(chan Job, queueSize),
        quit:     make(chan bool),
    }
}

func (w *Worker) Start() {
    go func() {
        for {
            select {
            case job := <-w.jobQueue:
                ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
                if err := job.Execute(ctx); err != nil {
                    log.Printf("Job failed: %v", err)
                }
                cancel()
            case <-w.quit:
                return
            }
        }
    }()
}

func (w *Worker) Enqueue(job Job) {
    w.jobQueue <- job
}

func (w *Worker) Stop() {
    w.quit <- true
}
```

## Next Steps

- Explore [Performance Optimization](./performance)
- Learn about [Scaling Strategies](./scaling)
- See [Production Deployment](../deployment/production)

# Chapter 3 — Chapter Summary

## Book: Разработка приложений в микросервисной архитектуре с нуля
## Author: Юлия Апапова

### Chapter Title
Способы взаимодействия между микросервисами (Ways of Communication Between Microservices)

### Chapter Pages
- Start: 110
- End: 157
- Total: 48 pages

### Chapter Goal
Introduce the patterns and technologies for microservice communication (HTTP, gRPC, message brokers) and implement the **API Gateway** — a central facade service that handles authentication, request routing, and data aggregation across internal services.

---

## Key Architecture Patterns

### 1. Synchronous Communication
- **HTTP**: Text-based protocol with request/response model. Used for external-facing APIs.
  - Methods: GET, POST, PUT, PATCH, DELETE
  - Status codes: 1xx (info), 2xx (success), 3xx (redirect), 4xx (client error), 5xx (server error)
  - Runs over TCP/IP using HTTP/1.1 or HTTP/2

- **gRPC**: Binary protocol using HTTP/2 + Protocol Buffers (Protobuf). Used for efficient service-to-service communication.
  - Service definitions in `.proto` files with `service` and `rpc` blocks
  - 4 communication patterns: Unary, Server-streaming, Client-streaming, Bidirectional streaming
  - Cross-language compatibility (Go, Java, C++, Python, etc.)
  - Testing via BloomRPC or Envoy (not Postman directly)

### 2. Asynchronous Communication
- **RabbitMQ (AMQP)**: Message broker with exchange/queue pattern
  - Exchange types: Fanout, Direct, Topic
  - Delivery guarantees: AtMostOnce, AtLeastOnce, ExactlyOnce
  - Channels within connections for concurrent operations
  - Used for background processing and event bus patterns

- **Kafka**: Distributed streaming platform with pull model
  - Topics → Partitions (ordered, indexed by offset)
  - Retains messages after consumption (unlike RabbitMQ)
  - Pull model: consumers poll for batches
  - Topic types: Regular (time/size based retention) and Compacted (key-based retention)
  - Kafka Connect for external system integration
  - Streams API for real-time processing

- **Redis**: In-memory key-value store
  - Complements disk-based databases (PostgreSQL/MySQL/MongoDB)
  - Use cases: caching, session storage, queues, geospatial data
  - Data structures: strings, hashes, lists, sets, sorted sets, bitmaps, HyperLogLog, streams
  - Not a full replacement for persistent DBs (expensive memory)

### 3. API Gateway Pattern
The **Gateway** service is the single entry point for all external clients:
- Validates authentication on every request via JWT Interceptor
- Routes requests to the correct internal service (Auth or Account)
- Transforms request/response formats (proto ↔ internal models)
- Aggregates data from multiple services (e.g., profile + photo + statistics)
- Centralizes cross-cutting concerns: token validation, rate limiting, caching, API versioning

#### Gateway Service Endpoints

| Method | gRPC Method | HTTP | Auth Required |
|--------|-------------|------|---------------|
| Register | Register | POST /api/v1/auth/register | No |
| Login | Login | POST /api/v1/auth/login | No |
| Refresh | Refresh | POST /api/v1/auth/refresh | No |
| Logout | Logout | POST /api/v1/auth/logout | No |
| Validate | ValidateToken | POST /api/v1/auth/validate | No |
| Create User | CreateUser | POST /api/v1/users | Yes |
| Get User | GetUser | GET /api/v1/users/{user_id} | Yes |
| Get Current User | GetCurrentUser | GET /api/v1/users/me | Yes |
| List Users | GetUsers | GET /api/v1/users | Yes |
| Update User | UpdateUser | PUT /api/v1/users/{user_id} | Yes |
| Update Me | UpdateCurrentUser | PUT /api/v1/users/me | Yes |
| Delete User | DeleteUser | DELETE /api/v1/users/{user_id} | Yes |
| Delete Me | DeleteCurrentUser | DELETE /api/v1/users/me | Yes |

### 4. JWT Interceptor Architecture

The **JWT Interceptor** centralizes authentication:

1. **Check if method is public** → skip authentication
2. **Extract Bearer token** from gRPC metadata (Authorization header)
3. **Validate JWT** → parse, verify signature (HMAC), check expiry
4. **Extract user_id** from JWT claims
5. **Add user_id to context** → downstream handlers receive authenticated user
6. **Forward to business logic**

#### Key Design Decisions:
- Custom `JWTClaims` struct embeds `jwt.RegisteredClaims` + `UserID` field
- Public methods (Register, Login, Refresh) bypass authentication
- Token validation happens once at Gateway; internal services trust context
- `ValidateTokenWithoutAuthService` enables local JWT validation (no round-trip to Auth service)
- Both Unary and Stream interceptors are implemented

---

## Service Architecture Summary

| Service | Port | Purpose | DB Port |
|---------|------|---------|---------|
| Account | 50051 | User profiles (CRUD) | 5432 |
| Auth | 50052 | Authentication (credentials, tokens) | 5433 |
| Gateway | 50053 | Unified entry point (facade) | N/A |

### Docker Compose Configuration
- Three PostgreSQL containers (one per service) on ports 5432, 5433, 5434
- Shared `app-network` (bridge driver)
- Health checks using `pg_isready`
- Named volumes per database

### Registration Flow (Cross-service)
```
Client → Gateway.Register → Account.CreateUser (creates user, NO password) → Auth.Register (stores bcrypt hash) → Auth.Login (returns token pair) → Client
```

### Request Flow (Authenticated)
```
Client → Gateway (JWT Interceptor validates token) → Gateway forwards to Account/Auth → Response returned through Gateway to Client
```

---

## Technologies

| Technology | Purpose |
|-----------|---------|
| Go | Language |
| gRPC | Service-to-service synchronous communication |
| Protobuf | gRPC message format (binary, cross-language) |
| HTTP | External API protocol |
| AMQP | RabbitMQ messaging protocol |
| RabbitMQ | Message broker (async, push model) |
| Apache Kafka | Streaming platform (async, pull model) |
| Redis | In-memory key-value store (cache/queue) |
| PostgreSQL | Persistent storage (one DB per service) |
| bcrypt | Password hashing (in Auth service) |
| JWT (HMAC-SHA256) | Access token format |
| pg_isready | PostgreSQL health check |
| goose/v3 | Database migrations |
| zerolog | Structured logging |

## What's Next

Chapter 4 covers the **Transaction service** — implementing money transfers between accounts using database transactions, ACID properties, isolation levels, and additional database migrations/indexes.

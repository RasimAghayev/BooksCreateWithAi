# Account Service

## Environment Variables
| Name | Description | Default |
|------|-------------|---------|
| `DB_DSN` | PostgreSQL connection string | - |
| `KAFKA_BROKER_HOST` | Kafka broker address | localhost:9092 |
| `KAFKA_CONSUMER_GROUP` | Kafka consumer group | consumer-account |
| `KAFKA_TRANSACTION_TOPIC` | Kafka topic for transaction requests | transaction_request |
| `GRPC_PORT` | gRPC server port | 50051 |

## Running Locally
```bash
go run ./cmd
```

## Running with Docker
```bash
docker build -t account:latest .
docker run -p 50051:50051 -e DB_DSN=postgres://postgres:postgres@host:5432/account?sslmode=disable account:latest
```

## Project Structure
- `internal/domain/` — domain models
- `internal/repository/` — GORM/SQL repository
- `internal/service/` — business logic
- `internal/transport/` — gRPC transport layer
- `internal/config/` — configuration
- `internal/migrations/` — database migrations
# Auth Service

## Environment Variables
| Name | Description | Default |
|------|-------------|---------|
| `DB_DSN` | PostgreSQL connection string | - |
| `JWT_SECRET` | JWT signing secret | secret |
| `ACCESS_TOKEN_TTL_MINUTES` | Access token lifetime in minutes | 60 |
| `REFRESH_TOKEN_TTL_DAYS` | Refresh token lifetime in days | 30 |
| `GRPC_PORT` | gRPC server port | 50052 |

## Running Locally
```bash
go run ./cmd
```

## Running with Docker
```bash
docker build -t auth:latest .
docker run -p 50052:50052 -e DB_DSN=postgres://postgres:postgres@host:5432/auth?sslmode=disable auth:latest
```

## Project Structure
- `internal/domain/` — domain models
- `internal/repository/` — SQL repository
- `internal/service/` — business logic (bcrypt, JWT)
- `internal/transport/` — gRPC transport layer
- `internal/config/` — configuration
- `internal/migrations/` — database migrations
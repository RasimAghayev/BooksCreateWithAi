# Gateway Service

API Gateway that routes HTTP requests to gRPC microservices.

## Environment Variables
| Name | Description | Default |
|------|-------------|---------|
| `ACCOUNT_GRPC_HOST` | Account service gRPC address | localhost:50051 |
| `AUTH_GRPC_HOST` | Auth service gRPC address | localhost:50052 |
| `TRANSACTION_GRPC_HOST` | Transaction service gRPC address | localhost:50054 |
| `JWT_SECRET` | JWT signing secret | secret |
| `HTTP_PORT` | HTTP server port | 8080 |

## Running Locally
```bash
go run ./cmd
```

## Running with Docker
```bash
docker build -t gateway:latest .
docker run -p 8080:8080 gateway:latest
```

## API Endpoints
- `POST /auth/register` — Register new user
- `POST /auth/login` — Login with email/password
- `POST /auth/refresh` — Refresh access token
- `POST /auth/logout` — Logout
- `GET /auth/me` — Get current user
- `POST /account/deposit` — Deposit to account
- `POST /account/withdraw` — Withdraw from account
- `POST /account/transfer` — Transfer between accounts
- `GET /account/balance/:account_id` — Get balance
- `POST /transaction/deposit` — Create deposit transaction
- `POST /transaction/withdraw` — Create withdraw transaction
- `POST /transaction/transfer` — Create transfer transaction
- `GET /transactions?user_id=1&limit=10&offset=0` — Get transactions
- `GET /health` — Health check
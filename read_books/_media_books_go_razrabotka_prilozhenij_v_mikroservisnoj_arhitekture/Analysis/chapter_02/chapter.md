# Chapter 2 — Chapter Summary

## Book: Разработка приложений в микросервисной архитектории с нуля
## Author: Юлия Апапова

### Chapter Title
Разработка микросервиса авторизации и аутентификации (Auth)

### Chapter Pages
- Start: 87
- End: 115
- Total: 29 pages

### Chapter Goal
Build a secure Auth service that handles the full authentication and authorization lifecycle for the microservice system. This service is separate from Account (User service) and manages user credentials, password hashing, JWT access tokens, and refresh token rotation.

### Key Architecture Decisions

1. **Separate Auth Service**: Auth is isolated from Account. Account manages user profiles; Auth manages credentials, sessions, and tokens.

2. **Two Token Strategy**:
   - **Access Token (JWT)**: Short-lived (60 min), stateless, contains user ID + expiry. Used by other services for request authorization.
   - **Refresh Token**: Long-lived (30 days), stateful (stored in DB as SHA-256 hash), revocable. Used to obtain new access tokens without re-login.

3. **Repository Pattern**: Database access is abstracted behind the `Repository` interface. The AuthService depends on this interface, enabling easy testing and DB swapping.

4. **Manual Dependency Injection**: Go DI without framework — AuthService receives `*config.Config`, `*zerolog.Logger`, and `Repository` via constructor.

5. **Security by Design**:
   - Passwords hashed with bcrypt (salt + work factor 10)
   - Refresh tokens hashed with SHA-256 before DB storage
   - JWT signed with HMAC-SHA256 using server-side secret
   - Soft-delete refresh tokens via `revoked_at` timestamp

### Data Model

```
users table:
  - id (BIGSERIAL PK)
  - login (TEXT UNIQUE)
  - email (TEXT UNIQUE)
  - password_hash (TEXT) — bcrypt hash
  - created_at, updated_at (TIMESTAMP)

refresh_tokens table:
  - id (BIGSERIAL PK)
  - user_id (BIGINT FK → users.id, ON DELETE CASCADE)
  - token (TEXT UNIQUE) — SHA-256 hash
  - expires_at (TIMESTAMP)
  - revoked_at (TIMESTAMP NULL) — soft revoke
  - created_at (TIMESTAMP)
```

### Service Interface

AuthService implements:
- `Register(ctx, RegisterRequest)` → hashes password, saves user
- `Login(ctx, LoginRequest)` → verifies password, issues token pair
- `Refresh(ctx, refreshToken)` → validates refresh token, issues new pair
- `Validate(ctx, accessToken)` → verifies JWT, returns user ID
- `Logout(ctx, refreshToken)` → revoke refresh token (soft delete)

### Token Lifecycle

```
Register → User saved with bcrypt password hash
  ↓
Login → bcrypt.CompareHashAndPassword → issueTokens(access + refresh)
  ↓            ↖
Access token used → Validate → user ID returned
  ↓ (expires)
Refresh → GetRefreshToken → check revoked/expiry → issueTokens(new pair)
  ↓
Logout → RevokeRefreshToken → revoked_at = NOW()
```

### Configuration

| Variable | Default | Purpose |
|----------|---------|---------|
| SERVICE_NAME | auth-service | Service identifier |
| GRPC_PORT | 50052 | gRPC server port |
| LOG_LEVEL | info | Logging verbosity |
| DB_DSN | (required) | Database connection string |
| JWT_SECRET | (required) | HMAC secret for signing access tokens |
| ACCESS_TOKEN_TTL_MINUTES | 60 | Access token lifetime |
| REFRESH_TOKEN_TTL_DAYS | 30 | Refresh token lifetime |

### Technologies Used

| Technology | Purpose |
|-----------|---------|
| Go | Language |
| gRPC + Protobuf | Service-to-service communication |
| goose/v3 | Database migrations |
| GORM | ORM |
| bcrypt | Password hashing |
| jwt/v5 | JWT creation/validation |
| zerolog | Structured logging |
| env/v10 | Environment variable config |
| godotenv | .env file support |

### What's Next

The Auth service needs:
1. gRPC server implementation (handler layer translating proto RPC to service calls)
2. Server startup boilerplate (main.go)
3. Inter-service communication setup (Auth ↔ Account)

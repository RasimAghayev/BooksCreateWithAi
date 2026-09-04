# Chapter 4 — Chapter Summary

## Book: Разработка приложений в микросервисной архитектуре с нуля
## Author: Юлия Апапова

### Chapter Title
Разработка модуля Transaction (Developing the Transaction Module)

### Chapter Pages
- Start: 158
- End: 264
- Total: ~107 pages (Chapter 4 + start of Chapter 5 references)

### Chapter Goal
Implement the **Transaction service** for handling monetary operations (deposit, withdraw, transfer) using **double-entry bookkeeping** with database-level atomicity (ACID), and integrate it with the **Account service** using the **Saga pattern** with Kafka for distributed transaction coordination. The chapter also covers advanced database theory (BCNF, 4NF, 5NF) and PostgreSQL indexing strategies.

---

## Key Architecture Decisions

### 1. Distributed Transaction Problem
- Transaction service and Account service use **separate databases** → can't use single ACID transaction across both
- Solution: **Saga pattern** with Kafka event broker — eventual consistency

### 2. Saga Pattern Implementation
- **Phase 1**: Transaction service creates transaction (pending) → publishes request to Kafka → Account service consumes and processes balance change → publishes response → Transaction service updates status (completed/failed)
- **Compensation**: If Kafka publish fails, transaction is marked as `failed` immediately
- **Idempotency**: `operation_id` (transaction ID) prevents duplicate processing
- **Topics**: `transaction_data` (requests), `transaction_response` (responses)

### 3. Double-Entry Bookkeeping
- Every Transaction has 1+ TransactionEntry records with DEBIT/CREDIT directions
- Deposit: 1 DEBIT entry (increase balance)
- Withdraw: 1 CREDIT entry (decrease balance)
- Transfer: 1 CREDIT (sender) + 1 DEBIT (recipient) — within a single DB transaction
- Amounts stored in **kopeyki** (int64) — `int64(amount * 100)` to avoid float precision errors

### 4. Transaction Lifecycle
```
1. CREATE Transaction (status=pending)
2. CREATE Entry/Entries (DEBIT and/or CREDIT)
3. [Saga] Publish to Kafka → Account updates balance
4. [Saga] Receive response → UPDATE Transaction status
   - completed (success) or failed (failure)
```

### 5. Account Service Extensions
- **New fields**: `balance` (float, default 0.00), `is_deleted` (bool, default false)
- **Soft delete**: `DeleteUser` sets `is_deleted=true` instead of hard delete — preserves transaction history
- **Balance operations**: `GetBalance`, `UpdateBalance` (with sufficient funds check), `TransferBalance` (atomic, within DB transaction)
- **Validation**: Negative amounts rejected; transfers to self rejected; insufficient funds checked before debiting

### 6. Database Design
- **Normal Forms**: BCNF (Boyce-Codd), 4NF, 5NF, Domain-Key NF, 6NF
- **Indexing**: B-tree indexes on `user_id`, `status`, `created_at` (transactions); `login` (account)
- **Docker Compose**: Three PostgreSQL containers (Account:5432, Auth:5433, Transaction:5434) + Kafka

---

## Service Architecture (Final State)

```
                    ┌─────────────┐
                    │ API Gateway │ (Port 50053)
                    │  (Facade)   │ JWT Interceptor
                    └──────┬──────┘
          ┌───────────────┼────────────────┐
          │ gRPC (50051)  │                │
          ▼               │                ▼
   ┌─────────────┐      │ gRPC (50052)   ┌──────────────┐
   │ Account Svc │      │                │    Auth Svc  │
   │             │      │                │              │
   │ + balance   │      │                │ + JWT tokens │
   │ + is_deleted│      │                │ + bcrypt     │
   │ + Kafka     │      │                │              │
   └──────┬──────┘      │                └──────────────┘
          │             │                         ▲
          │ Kafka ◄─────┼──────── Kafka ──────── │
          │             │                         │
          ▼             │                         │
   ┌─────────────┐      │                        │
   │ Transaction │      │                        │
   │   Service   │ ──────┼───────────────────────┘
   │   (50054)   │ gRPC
   │             │
   │ + Transaction│
   │ + Entries    │
   │ + Saga       │
   │ + Kafka      │
   └─────────────┘
```

---

## Microservices Inventory

| Service | Port | gRPC Port | DB | DB Port | Purpose |
|---------|------|-----------|----|---------|---------|
| Account | 50051 | 50051 | PostgreSQL | 5432 | User profiles + balance management |
| Auth | 50052 | 50052 | PostgreSQL | 5433 | Authentication (bcrypt, JWT, refresh tokens) |
| Gateway | 50053 | 50053 | N/A | N/A | API facade (JWT interception, routing, aggregation) |
| Transaction | 50054 | 50054 | PostgreSQL | 5434 | Financial transactions (deposit/withdraw/transfer) |

### Docker Services
| Service | Port | Purpose |
|---------|------|---------|
| Kafka | 29092/9092 | Message broker (Saga pattern) |
| Kafka UI | 29093:8080 | Kafka monitoring UI |

---

## Transaction States

| State | Description |
|-------|-------------|
| `pending` | Transaction created, awaiting balance operation result |
| `completed` | Balance operation succeeded |
| `failed` | Balance operation failed or Kafka error |
| `cancelled` | Manual cancellation (not yet implemented) |

## Transaction Types

| Type | DEBIT Entry | CREDIT Entry |
|------|-------------|--------------|
| Deposit | ✓ (increase) | ✗ |
| Withdraw | ✗ | ✓ (decrease) |
| Transfer | ✓ (recipient) | ✓ (sender) |

## Balance Operation Types

| Type | Effect |
|------|--------|
| DEPOSIT | `balance += amount` |
| CREDIT | `balance -= amount` (with sufficient funds check) |

---

## Key Technologies

| Technology | Purpose |
|-----------|---------|
| Go | Service implementation language |
| gRPC + Protobuf | Service-to-service communication |
| Kafka (confluentinc/cp-kafka) | Async event broker (Saga pattern) |
| segmentio/kafka-go | Go Kafka client library |
| PostgreSQL 15 | Persistent storage (4 instances) |
| GORM | ORM with transaction support |
| goose/v3 | Database migration tool |
| bcrypt | Password hashing (Auth service) |
| JWT (HMAC-SHA256) | Access token format |
| Docker Compose | Local infrastructure orchestration |

## What's Next

Chapter 5 covers **microservice deployment** — strategies (single VM, multi-VM, containers, Kubernetes, serverless) and practical nginx configuration examples. The complete system is built but deployment to production is not yet covered.

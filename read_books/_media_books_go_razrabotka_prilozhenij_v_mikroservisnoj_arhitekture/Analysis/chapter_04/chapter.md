# Chapter 4 — Chapter Summary

## Book: Разработка приложений в микросервисной архитектуре с нуля
## Author: Юлия Апапова

### Chapter Title
Разработка модуля Transaction (Developing the Transaction Module)

### Chapter Pages
- Start: 158
- End: 204
- Total: 47 pages

### Chapter Goal
Build the **Transaction service** for handling money transfers (deposit, withdraw, transfer) using double-entry bookkeeping with database-level atomicity. Also extends the Account service with balance management and soft-delete functionality.

---

## Key Concepts

### Database Design Theory (Pages 158–165)
1. **Normal Forms** (beyond 3NF covered in Chapter 1):
   - **НФБК (Boyce-Codd)**: Eliminates redundancy when tables have multiple composite candidate keys
   - **4NF**: Eliminates multivalued dependencies
   - **5NF**: Eliminates join dependencies
   - **ДКНФ**: All constraints are consequences of domain + key constraints
   - **6NF**: All join dependencies preserved (rarely used in practice)

2. **ACID Properties**:
   - **Atomicity**: All-or-nothing (rollback on error)
   - **Consistency**: Valid state before and after
   - **Isolation**: Concurrent transactions don't interfere
   - **Durability**: Committed changes persist

3. **Transaction Isolation Levels**:
   - Read Uncommitted → Read Committed → Repeatable Read → Serializable
   - Concurrency anomalies: lost update, dirty read, non-repeatable read, phantom read, serialization anomaly

### Transaction Service (Pages 167–186)
- **Service purpose**: Financial operations (deposit, withdraw, transfer)
- **Data model**: Transaction (header) + TransactionEntry (line items, DEBIT/CREDIT)
- **Double-entry bookkeeping**: Every credit has matching debit
- **Amount storage**: `int64` in kopeyki/cents (`amount * 100`) to avoid float precision
- **Transaction lifecycle**: CREATE (pending) → ADD entries → UPDATE (completed)
- **GORM Transaction()**: Wraps all operations for atomicity
- **No delete**: Financial transactions are immutable (audit requirement)
- **Repository interface**: Defined in service layer (Go DI)

### Account Service Integration (Pages 187–204)
- **New Account fields**: `balance` (float), `is_deleted` (bool, soft delete)
- **Soft delete pattern**: `is_deleted` flag instead of hard DELETE to preserve referential integrity
- **Balance methods**: `GetBalance`, `UpdateBalance` (deposit/credit with validation), `TransferBalance` (atomic)
- **Balance validation**: Check `newBalance >= 0` before crediting (prevents negative balance)
- **Transfer atomicity**: Both debit + credit in single DB transaction
- **Service layer**: Thin wrapper with input validation + structured logging
- **Proto contract**: New RPCs for Deposit, Withdraw, Transfer, GetBalance with operation_id for idempotency

### Infrastructure (Pages 167)
- **Docker Compose**: Three PostgreSQL containers (Account:5432, Auth:5433, Transaction:5434)
- **Indexes**: Added on user_id, status, created_at for query performance
- **Migration pattern**: Named migrations with up/down functions, `goose.AddNamedMigrationContext()`

---

## Architecture Summary

```
                    ┌─────────────┐
                    │ API Gateway │ (Port 50053)
                    │  (Facade)   │
                    └──────┬──────┘
         Auth (gRPC) ┌─────┴─────┐
                     │ Auth Svc  │ (Port 50052)
                     └──────────┘
        Balance gRPC ┌──────────────┐
                     │ Account Svc  │ (Port 50051)
                     └──────┬───────┘
                            │ DB Trans
                            │
                    ┌───────┴─────┐
                    │ Transaction │ (Balance updates)
                    │   Svc       │ Port: N/A (calls Account gRPC)
                    └─────────────┘
```

### Microservice Communication Flow

1. **Client → Gateway**: Single entry point, JWT validated by interceptor
2. **Gateway → Auth**: Token issuance (gRPC)
3. **Gateway → Account**: User management (gRPC)
4. **Transaction → Account**: Balance operations via gRPC calls

---

## Data Models

### Transaction Service
```
transactions:        transaction_entries:
  id (PK)              id (PK)
  user_id               transaction_id (FK → transactions)
  amount (int64, cents) account_id
  type (deposit/withdraw/transfer) direction (DEBIT/CREDIT)
  status (pending/completed/failed/cancelled) amount (decimal)
  created_at, updated_at
```

### Account Service (extended)
```
users:
  id (PK)
  login, email, phone, first_name, last_name, middle_name, age
  balance (decimal 15,2)     -- NEW
  is_deleted (boolean, default false)  -- NEW (soft delete)
  created_at, updated_at
```

---

## Transaction Flow (Transfer Example)

```
1. Gateway.Receive Transfer request
2. Gateway calls TransactionService.Transfer(user_id, recipient_id, amount)
3. Repository.Transfer():
   a. START DB TRANSACTION
   b. Create Transaction record (status=pending)
   c. Create DEBIT entry (sender loses funds)
   d. Create DEBIT entry (recipient gains funds)
   e. Update Account balances (via gRPC or same DB)
   f. UPDATE Transaction (status=completed)
   g. COMMIT/ROLLBACK
4. Return TransactionDetails with entries
```

---

## Technologies

| Technology | Purpose |
|-----------|---------|
| Go | Language |
| PostgreSQL | Persistent storage (3 instances) |
| gRPC + Protobuf | Service-to-service sync comm |
| GORM | ORM |
| goose/v3 | Database migrations |
| zerolog | Structured logging |
| Docker Compose | Local development environment |

## What's Next

Chapter 5 covers the final service implementation, integration, and potentially deployment/CI/CD setup.

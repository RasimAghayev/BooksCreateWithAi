# Chapter 4 — Разработка модуля Transaction

Pages: 157-266

## Bu chapter nədən bəhs edir?
Bu chapter Transaction (tranzaksiya) mikroservisini hazırlamağı və
paylanmış tranzaksiya problemlərini həll etməyi izah edir. Verilənlər
bazasının normallaşdırılması, miqrasiyalar, ACID prinsipləri və Saga
pattern-i əhatə edir.

## Əsas fikirlər
### 1. Verilənlər bazası layihələndirilməsi
- 1NF, 2NF, 3NF normal formaları
- Boyce-Codd Normal Form (BCNF / НФБК)
- 4NF, 5NF, 6NF
- Domain-Key Normal Form (DKNF)

### 2. Migrasiyalar və tranzaksiyalar
- Database migration mexanizmi
- ACID prinsipləri (Atomicity, Consistency, Isolation, Durability)
- Paralel tranzaksiyalar
- SQL-də izolasiya səviyyələri

### 3. Paylanmış tranzaksiyalar
- Two-Phase Commit (2PC / İki mərhəmətli fiqaslama)
- Saga pattern — compensasiya mexanizmi
- Orchestration vs Choreography

### 4. İnteqrasiya
- Transaction və Account mikroservislərinin inteqrasiyası
- Event-driven arxitektura

## Əsas terminlər
- Transaction (tranzaksiya)
- ACID (Atomicity, Consistency, Isolation, Durability)
- Migration (verilənlər bazası miqrasiyası)
- Index (indeks)
- Distributed Transaction (paylanmış tranzaksiya)
- Two-Phase Commit (iki mərhəmətli fiqaslama)
- Saga (kompensasiya pattern-i)
- Compensation (kompensasiya əməliyyatı)
- Event-Driven Architecture (hadisə əsaslı arxitektura)
- Idempotency (idempotenslik / təkrar emsallıq)

## Praktik nəticə
Transaction mikroservisi hazırlanır və Account ilə inteqrasiya olunur.
Paylanmış tranzaksiya problemləri Saga pattern-i ilə həll olunur.

## Mənbə
Pages: 157-266

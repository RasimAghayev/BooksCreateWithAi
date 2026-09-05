# 13. Understanding GraphQL with a Real Project

**Səhifələr:** 281-353

## Bu fəsil nədən bəhs edir?

Bu fəsil real bir layihə — login sistemi — qurur. Backend (Apollo Server + PostgreSQL + Sequelize + JWT) və Frontend (Apollo Client + React) hissələri tam şəkildə implementasiya olunur.

## Əsas fikirlər

### 1. GraphQL-in üstünlükləri

REST API-lardan fərqli olaraq GraphQL-də istifadəçi yalnız lazım olan sahələri sorğu edə bilər. Bir endpoint, çoxlu sorğu. Over-fetching və under-fetching problemləri həll olunur. Type sistemi schema-first-dir — sənədləşdirmə avtomatik.

### 2. Backend qurma

PostgreSQL quraşdır, `.env` faylı ilə konfiqurasiyaları saxla, Apollo Server qur, GraphQL types/queries/mutations/resolvers yaz, Sequelize ORM ilə User model yarat, JWT ilə authentication əlavə et.

### 3. Scalar types, Queries, Mutations

Scalar types — GraphQL-də mövcud olan əsas tiplər (Int, String, Boolean, ID). `Query` — məlumat oxuma (`getUsers`, `getUser(id)`). `Mutation` — məlumat dəyişdirmə (`createUser`, `login`).

### 4. Resolvers

Hər query/mutation üçün resolver funksiyası yazılır. Resolver GraphQL sorğusunu həll edib data qaytarır. Sequelize modelləri ilə əlaqə: `User.findOne({ where: { email } })`.

### 5. JWT authentication

`jsonwebtoken` paketi ilə token yarat. `sign({ userId }, secret, { expiresIn: '1h' })` istifadəçi üçün token yaradır. Sonra hər protected mutation-da token verify et.

### 6. Frontend qurma (Apollo Client)

Webpack 5 ilə frontend build qur, TypeScript konfiqurasiya et, Express server ilə serve et, Apollo Client yarat, GraphQL sorğularını komponentlərdən çağır, login flow qur, dashboard komponentləri yarat.

### 7. Testing login sistemi

Backend GraphQL sorgularını GraphQL Playground ilə test et. Model validation-ları test et. Frontend-də tam login axınını (daxil ol, çıxış et, token saxla) test et.

## Əsas terminlər

- GraphQL
- Apollo Server
- Apollo Client
- PostgreSQL
- Sequelize
- ORM
- JWT
- JSON Web Token
- Resolver
- Scalar Type
- Query
- Mutation
- Schema
- Type Definition

## Praktik nəticə

Bu fəsildəki tam login sistemini sıfırdan qur: PostgreSQL qur, backend və frontend yaz, login/ register/ dashboard səhifələri əlavə et, JWT ilə auth qoru.

## Mənbə

Pages: 281-353
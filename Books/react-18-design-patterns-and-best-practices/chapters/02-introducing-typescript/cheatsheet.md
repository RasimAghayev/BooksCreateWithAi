# Cheat Sheet — Introducing TypeScript

## TypeScript əsasları

### Type annotation

**Nə edir:** Dəyişənlərə tip vermək. Compile-time xəta tutmağa kömək edir.

**Kod:**

```ts
const name: string = 'Carlos';
const age: number = 30;
const active: boolean = true;
```

**Mənbə:** Chapter 2, page 44-57

### Interface

**Nə edir:** Obyektin formasını təyin edir. `?` optional field bildirir.

**Kod:**

```ts
interface User {
  id: number;
  name: string;
  email?: string;
}
const user: User = { id: 1, name: 'Ali' };
```

**Mənbə:** Chapter 2, page 44-57

### Enum

**Nə edir:** Sabit dəyərlər toplusu, type-safe alternativlər string literal-lara.

**Kod:**

```ts
enum Status {
  Active = 'ACTIVE',
  Inactive = 'INACTIVE',
}
```

**Mənbə:** Chapter 2, page 44-57

### Generics

**Nə edir:** Funksiyanı tip-agnostik saxlamaq üçün generic istifadə olunur.

**Kod:**

```ts
function identity<T>(arg: T): T {
  return arg;
}
const out = identity<string>('hello');
```

**Mənbə:** Chapter 2, page 44-57

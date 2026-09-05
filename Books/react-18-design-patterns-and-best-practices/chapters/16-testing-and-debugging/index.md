# 16. Testing and Debugging

**Səhifələr:** 421-466

## Bu fəsil nədən bəhs edir?

Bu fəsil testing-in əhəmiyyətini, Jest, React Testing Library və Vitest ilə test yazmağı, həmçinin React DevTools və Redux DevTools istifadəsini öyrədir.

## Əsas fikirlər

### 1. Testing-in faydaları

Testlər refactoring zamanı kömək edir, bug-ları erkən tutur, kodun davranışını sənədləşdirir, komanda arasında güvən yaradır. Test piramidi: unit (ən çox), integration, e2e (ən az).

### 2. Jest ilə JavaScript testing

Jest ən populyar test runner-dir. `describe`, `it/test`, `expect` ilə test yazılır. Mock, spy, snapshot xüsusiyyətləri var. React komponentləri üçün `@testing-library/react` ilə birlikdə istifadə olunur. `render`, `screen`, `fireEvent`/`userEvent` API-ləri ilə komponentləri test etmək asandır.

### 3. Event testing

`userEvent.click(button)` kimi real istifadəçi davranışını simulyasiya edir. `fireEvent` daha aşağı səviyyəli API-dir. Form göndərmək, input dəyişdirmək kimi ssenariləri test etmək olur.

### 4. Vitest — Vite ilə testing

Vitest Vite-ə inteqrasiya olunmuş test runner-dir — Jest API-sinə uyğundur, lakin daha sürətlidir (ESM, native ES modules). `globals: true` ilə `describe`, `it` import etmədən istifadə oluna bilər. In-source testing dəstəyi var (`if (import.meta.vitest)` bloku).

### 5. React DevTools

Brauzer extension-ı komponent ağacını, props, state, hooks göstərir. Profiler tab render vaxtlarını ölçür. Production-da istifadə üçün `__REACT_DEVTOOLS_GLOBAL_HOOK__` istifadə olunur.

### 6. Redux DevTools

Redux store-un vəziyyətini, action-ları, time-travel debugging imkan verir. `redux-devtools-extension` middleware ilə qoşulur.

## Əsas terminlər

- Testing
- Jest
- Vitest
- React Testing Library
- describe
- it
- expect
- Mock
- Snapshot
- userEvent
- fireEvent
- Profiler
- React DevTools
- Redux DevTools

## Praktik nəticə

Bir komponent üçün unit testlər yaz (render, klik hadisəsi, state dəyişikliyi). Sonra Vitest-ə migrate et və sürət fərqini ölç.

## Mənbə

Pages: 421-466
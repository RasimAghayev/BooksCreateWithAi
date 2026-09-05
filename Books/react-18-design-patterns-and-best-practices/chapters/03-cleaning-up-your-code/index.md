# 3. Cleaning Up Your Code

**Səhifələr:** 58-80

## Bu fəsil nədən bəhs edir?

Bu fəsil JSX-in əsaslarını, Prettier və ESLint ilə kod formatlama/qoruma qaydalarını, və funksional proqramlaşdırma prinsiplərini öyrədir.

## Əsas fikirlər

### 1. JSX nədir

JSX (JavaScript XML) JavaScript daxilində HTML-ə oxşar sintaksisdir. Babel kimi alətlər JSX-i React.createElement çağırışlarına çevirir. Bu, məntiq və UI-ı eyni yerdə saxlamağa imkan verir.

### 2. JSX-in HTML-dən fərqləri

`class` əvəzinə `className`, `for` əvəzinə `htmlFor`. Bütün atributlar camelCase ilə yazılır (`onclick` → `onClick`). Boş teqlər `<br />` kimi yazılmalıdır. Boolean atributlar (məsələn, `disabled`) `={true}` ilə verilir.

### 3. Babel və transpilasiya

Babel müasir JavaScript/TypeScript/JSX-i köhnə brauzerlərə uyğun JavaScript-ə çevirir. `@babel/preset-env`, `@babel/preset-react`, `@babel/preset-typescript` ən çox istifadə olunan preset-lərdir.

### 4. Prettier və ESLint

Prettier avtomatik formatlayıcıdır (boşluq, sətir uzunluğu). ESLint isə kod keyfiyyətini yoxlayır (sintaksis xətaları, anti-pattern-lər). `@typescript-eslint` paketi TypeScript üçün qaydalar əlavə edir. Git Hooks (husky, lint-staged) ilə commit zamanı avtomatik yoxlama qurulur.

### 5. Funksional proqramlaşdırma (FP) əsasları

FP-nin React üçün əhəmiyyəti: (1) First-class functions — funksiyalar dəyər kimi ötürülə bilər; (2) Purity — eyni giriş həmişə eyni çıxışı verir, yan təsirsiz; (3) Immutability — dəyərlər yaradıldıqdan sonra dəyişdirilmir; (4) Currying — çox arqumentli funksiyanı tək arqumentli funksiyalara bölmək; (5) Composition — kiçik funksiyaları birləşdirib böyük funksiya qurmaq. Bu prinsiplər React Hooks-un əsasını təşkil edir.

## Əsas terminlər

- JSX
- Babel
- Transpilation
- ESLint
- Prettier
- EditorConfig
- Functional Programming
- Purity
- Immutability
- Currying
- Composition

## Praktik nəticə

Layihənə ESLint + Prettier qur, sonra bir neçə komponenti FP prinsiplərinə uyğun yaz: pure funksiyalar, immutable state yeniləmələri, composition pattern.

## Mənbə

Pages: 58-80
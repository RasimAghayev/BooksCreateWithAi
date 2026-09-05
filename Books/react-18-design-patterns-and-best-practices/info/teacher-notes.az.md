# Müəllim Qeydləri — React 18 Design Patterns and Best Practices

**Müəllif:** Carlos Santana Roldán  
**İl:** 2023

## 📖 Kitab deyir:

Bu kitab dördüncü nəşrdir. React 18-in bütün yeni xüsusiyyətlərini (concurrent mode, automatic batching, transitions, useId, useTransition, useDeferredValue, useInsertionEffect) praktik nümunələrlə öyrədir. Kitab həm də TypeScript, Vite, React Router v6.4, Apollo GraphQL, SWR, Redux Toolkit, Next.js, Vitest və CircleCI kimi müasir ekosistem alətlərini əhatə edir.

1. **Fundamentals:** Declarative proqramlaşdırma, JSX, komponent arxitekturası
2. **TypeScript:** Statik tipləmə React layihələrində
3. **Patterns:** Composition, HOC, Container/Presentational, FunctionAsChild
4. **Hooks:** useState, useEffect, useMemo, useCallback, useReducer
5. **Data management:** Context API, SWR, Redux Toolkit
6. **SSR:** Universal tətbiqlər və Next.js
7. **Real project:** GraphQL login sistemi (backend + frontend)
8. **MonoRepo:** NPM Workspaces, Webpack, shared packages
9. **Performance və Testing:** Reconciliation, key'lar, Jest, Vitest
10. **Deployment:** DigitalOcean, nginx, PM2, CircleCI

## 👨‍🏫 Müəllim qeydi:

Bu kitab real layihə öyrənmək istəyən developerlər üçün idealdır. Chapter 13-də GraphQL login sistemi qurulur — bu, tələbələrə backend və frontend arasındakı əlaqəni göstərir. Chapter 14-də MonoRepo arxitekturası real şirkət ssenarisi kimi təqdim olunur. Bu iki fəsil birlikdə "production-ready" layihənin bütün aspektlərini göstərir.

### Ən vacib 5 fikir

1. **Composition > Inheritance:** React-də inheritance yoxdur. Komponentləri composition vasitəsilə birləşdirin. `children` prop və HOC pattern bu yanaşmanın əsasını təşkil edir.
2. **Hooks qaydaları:** Hooks yalnız React funksiyalarında və top-level-də çağırılmalıdır. Bu, state-in düzgün saxlanmasını təmin edir və `react-hooks/rules-of-hooks` ESLint qaydası ilə qorunur.
3. **Immutability vacibdir:** State və props həmişə yeni obyektlər kimi dəyişdirilməlidir. Bu, React.memo və useMemo'nun düzgün işləməsini təmin edir.
4. **Key prop'un əhəmiyyəti:** Siyahılarda index əvəzinə unikal ID istifadə edin. Index istifadəsi reconciliation zamanı performans problemlərinə və buglara səbəb olur.
5. **Memoization hər yerdə lazım deyil:** useMemo, useCallback, memo performans üçündür, lakin hər yerdə istifadə olunmaları əks effekt verə bilər (öz-özünə yaddaş xərcləri).

### Kitabın ən dəyərli hissəsi

**Chapter 13 — Understanding GraphQL with a Real Project (səhifələr 281-353)** — çünki bu fəsil real bir layihənin (login sistemi) həm backend (Apollo Server + Sequelize + PostgreSQL + JWT), həm frontend (Apollo Client + React) hissələrini praktik şəkildə göstərir. Tələbələr bu fəsildən sonra GraphQL-i production layihələrdə tətbiq edə bilərlər.

### Tövsiyələr

- **Praktika:** Hər chapter-ı oxuduqdan sonra Vite ilə kiçik bir React 18 layihəsi qurun.
- **TypeScript:** Əgər TypeScript ilə işləmirsinizsə, əvvəlcə Chapter 2-ni dərindən oxuyun, sonra digər fəsillərə keçin.
- **Hooks dərindən:** Chapter 8-i (React Hooks) diqqətlə oxuyun. Bu, müasir React-in əsasını təşkil edir.
- **GraphQL:** Chapter 13-ü oxuyarkən Apollo Server sənədlərini də paralel olaraq açın.
- **MonoRepo:** NPM Workspaces haqqında əlavə mənbələrdən oxuyun (Nx, Turborepo alternativləri haqqında da məlumatlanın).
- **Testing:** Chapter 16 Vitest-i Jest-ə nisbətən daha çox töv edir. Vite layihələri üçün Vitest daha sürətlidir.
- **Deployment:** Chapter 17 DigitalOcean istifadə edir, lakin AWS, Vercel və ya Netlify kimi alternativlər də öyrənilməlidir.

### Çətinliklər

- Chapter 14 (MonoRepo) çox uzundur (59 səhifə) və Webpack konfiqurasiyası detallıdır. Tələbələr burada yorula bilər — bu fəsli 2-3 hissəyə bölmək məsləhətdir.
- Chapter 13 (GraphQL) də uzundur (73 səhifə) — amma real layihə üçün bu detallı yanaşma vacibdir.

### Əlavə mənbələr

- React rəsmi sənədləri: https://react.dev
- TypeScript sənədləri: https://www.typescriptlang.org
- Apollo GraphQL sənədləri: https://www.apollographql.com/docs
- React Router sənədləri: https://reactrouter.com
- Vitest sənədləri: https://vitest.dev
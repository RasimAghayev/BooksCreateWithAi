# 12. Server-Side Rendering

**Səhifələr:** 259-280

## Bu fəsil nədən bəhs edir?

Server-Side Rendering (SSR) — HTML-in server tərəfində generasiya olunması. SEO, performans və social sharing üçün vacibdir. Bu fəsil əl ilə SSR qurmağı və Next.js istifadəsini izah edir.

## Əsas fikirlər

### 1. Universal application anlayışı

Universal (və ya isomorphic) tətbiq həm server, həm client tərəfində eyni kodla işləyir. React komponentləri server-də render olunur, sonra client-də hydrate olunur. Bu, bizə SEO-friendly və sürətli ilk yükləmə verir.

### 2. SSR-in səbəbləri

(1) SEO — Googlebot JS-rendered kontenti indeksləyə bilmir, lakin SSR kontentini görür. (2) Social sharing — Facebook, Twitter meta-tag-ləri oxuyur, SSR olmadan OG image düzgün görünmür. (3) Performance — istifadəçi HTML-i dərhal görür, JS yüklənməsini gözləmir (TTI azalır).

### 3. SSR çətinlikləri

Mürəkkəblik yüksəkdir: webpack server konfiqurasiyası, data fetching həm server həm client-də, hydration uyğunsuzluqları, server mühitinə aid API-lər (window, document) yoxdur. Next.js kimi framework-lər bu çətinlikləri gizlədir.

### 4. Əl ilə SSR qurma

Webpack ilə server bundle yarat, Express ilə sorğuları qəbul et, `renderToString` ilə HTML-i generasiya et, sonra client-ə göndər. Data fetching üçün `fetch` API istifadə olunur. `hydrateRoot` ilə client-də hydrate et.

### 5. Next.js ilə SSR

Next.js SSR-i default olaraq dəstəkləyir. Fayl sistemi əsaslı routing: `pages/index.tsx` → `/`. `getServerSideProps` hər sorğuda server-də çağırılır. Static Generation (`getStaticProps`) build zamanı HTML yaradır. Next.js 13+ App Router ilə React Server Components dəstəklənir.

## Əsas terminlər

- SSR
- Server-Side Rendering
- Universal App
- Isomorphic
- SEO
- Hydration
- renderToString
- Next.js
- getServerSideProps
- getStaticProps
- React Server Components
- App Router

## Praktik nəticə

Next.js ilə sadə bir blog qur: 3 səhifə, server-dən data fetch, SEO meta-tag-lər, deploy üçün hazır.

## Mənbə

Pages: 259-280
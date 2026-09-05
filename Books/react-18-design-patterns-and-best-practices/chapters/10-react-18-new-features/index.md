# 10. React 18 New Features

**Səhifələr:** 217-234

## Bu fəsil nədən bəhs edir?

React 18 köklü yeniliklər gətirdi: concurrent rendering, automatic batching, transitions, server-side Suspense, yeni API-lər (createRoot, hydrateRoot, renderToPipeableStream) və yeni hook-lar (useId, useTransition, useDeferredValue, useInsertionEffect).

## Əsas fikirlər

### 1. Concurrent mode

React 18 köhnə 'blocking' render modelini 'concurrent' (kəsilə bilən) modeli ilə əvəz etdi. İndi React böyük render işini dayandırıb daha vacib yeniləmələrə (məsələn, istifadəçi klikinə) cavab verə bilir. Bu, daha responsive UI deməkdir.

### 2. Automatic batching

React 17-də yalnız event handler-lər daxilində state yeniləmələri batch olunurdu. React 18-də bütün yerlərdə (Promise, setTimeout, native event handler-lər) avtomatik batch olunur. Nəticə: daha az render, daha yaxşı performans.

### 3. Transitions

`useTransition` hook-u uzunmüddətli state yeniləmələri üçün 'aşağı prioritetli' kimi işarələnməyə imkan verir. `startTransition` daxilindəki yeniləmələr daha az vacibdir və istifadəçi interaktivlığı kəsə bilər. `isPending` flag yüklənmə indikatoru üçün istifadə olunur.

### 4. Suspense on the server

`renderToPipeableStream` ilə server tərəfində HTML hissə-hissə göndərilə bilər. Suspense boundary-ləri hazır olan hissələri dərhal göndərir, yüklənməmiş hissələr sonra flush olunur. Nəticə: istifadəçi daha tez kontent görür (TTI azalır).

### 5. useId

`useId()` SSR-uyğun unikal ID generasiya edir. HTML-də eyni komponent server və client-də eyni ID almalıdır, əks halda hydration xətası yaranır. useId bu problemi həll edir.

### 6. useDeferredValue

`useDeferredValue(value)` value-nu prioritet az olan yeniləmə kimi işarələyir. Məsələn, axtarış input-u dəyəri dəyişəndə siyahı `useDeferredValue` ilə yavaş yenilənir ki, input responsiv qalsın.

### 7. useInsertionEffect

CSS-in-JS kitabxanaları (styled-components) üçün xüsusi hook. DOM mutation-dan əvvəl işləyir, CSS qaydalarını inject edir. `useEffect`-dən daha erkən çağırılır.

### 8. Strict Mode

React 18 Strict Mode hər komponenti iki dəfə render edir (development-da) — bu, side effect-lərin pure olub-olmadığını aşkar etmək üçündür. Production-da heç bir təsiri yoxdur.

## Əsas terminlər

- Concurrent Mode
- Batching
- Transition
- useTransition
- useDeferredValue
- useId
- useInsertionEffect
- createRoot
- hydrateRoot
- renderToPipeableStream
- Strict Mode
- Suspense
- Node.js 18
- LTS

## Praktik nəticə

Mövcud React 17 layihəsini React 18-ə upgrade et. Sonra böyük bir siyahını `useDeferredValue` ilə optimallaşdır.

## Mənbə

Pages: 217-234
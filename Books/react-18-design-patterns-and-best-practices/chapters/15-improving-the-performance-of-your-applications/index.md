# 15. Improving the Performance of Your Applications

**Səhifələr:** 413-420

## Bu fəsil nədən bəhs edir?

Bu fəsil React tətbiqlərinin performansını artırmaq üçün istifadə olunan texnika və alətləri qısaca izah edir.

## Əsas fikirlər

### 1. Reconciliation (Uzlaşdırma) necə işləyir

React hər render-də yeni Virtual DOM yaradır, köhnə ilə müqayisə edir (diff), fərqi real DOM-a tətbiq edir. Bu 'diffing' alqoritmi O(n) mürəkkəblikdədir, lakin bəzi hallarda (məs. siyahılarda key olmadan) daha aşağı performans ola bilər.

### 2. Keys-in düzgün istifadəsi

Siyahılarda hər elementə unikal, sabit ID key verilməlidir. Bu, React-ə hansı elementin hərəkət etdiyini, hansının silindiğini/əlavə olunduğunu bildirir. Index key istifadəsi performans problemlərinə səbəb olur.

### 3. Optimallaşdırma texnikaları

(1) Lazy loading — `React.lazy()` ilə komponentləri tələb üzrə yüklə. (2) Code splitting — Webpack ilə bundle-ları bölmə. (3) memo, useMemo, useCallback — lazım olan yerlərdə memoizasiya. (4) Virtualization — `react-window` kimi kitabxanalarla böyük siyahıları pəncərələmək.

### 4. Immutability vacibdir

State və props həmişə yeni obyektlər kimi 'dəyişdirilməlidir' (spread operatoru, slice, map, filter). Bu, `shallowCompare` və `React.memo`-nun düzgün işləməsini təmin edir. Mutasiya referens eyniliyini qoruyur və React 'dəyişiklik yoxdur' hesab edir.

### 5. Alətlər və kitabxanalar

Immutability üçün Immer (sadə sintaksis ilə immutable state yaratmaq). Babel plugins (lodash → tək funksiyalar, antd icons). Profiling alətləri: React DevTools Profiler, Chrome DevTools Performance tab.

## Əsas terminlər

- Reconciliation
- Diffing
- Virtual DOM
- Key Prop
- React.memo
- useMemo
- useCallback
- Lazy Loading
- Code Splitting
- Virtualization
- Immutability
- Immer
- Babel Plugin
- Profiler

## Praktik nəticə

React DevTools Profiler ilə tətbiqini yavaş render-ləri tap. Sonra siyahı komponentini virtualization ilə optimallaşdır (react-window istifadə edərək).

## Mənbə

Pages: 413-420
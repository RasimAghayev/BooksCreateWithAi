# 11. Managing Data

**Səhifələr:** 235-258

## Bu fəsil nədən bəhs edir?

React-də data idarəsinin 3 əsas yolu: Context API (lokal state-lər üçün), SWR (remote data + caching), Redux Toolkit (global state + complex business logic).

## Əsas fikirlər

### 1. Context API

`createContext` ilə context yaradılır, `<Provider value={...}>` ilə komponentlərə ötürülür, `useContext` ilə oxunur. Kiçik/middle layihələrdə props drilling həll edir. Amma bütün tətbiqi Context ilə idarə etmək performans problemi yaradır (bütün consumer-lər yenidən render olunur).

### 2. SWR (Stale-While-Revalidate)

Vercel-in data fetching kitabxanasıdır. Cache strategiyaları: stale-while-revalidate — dərhal cache-dən qaytarır, arxa planda yeniləyir. `useSWR(key, fetcher)` ilə sadə istifadə. Suspense ilə inteqrasiya mövcuddur.

### 3. React Suspense ilə SWR

Suspense məntiqi: data yüklənənə qədər fallback göstərilir. SWR `<Suspense fallback={<Skeleton />}>` ilə istifadə olunur. Pokedex nümunəsi ilə göstərilir: 151 Pokemon sorğusu paralel, skeleton loading, sonra data.

### 4. Redux Toolkit

Redux-un müasir, sadələşdirilmiş versiyasıdır. `createSlice` ilə reducer + action-lar bir yerdə yaradılır. `configureStore` ilə store qurulur. `useSelector` və `useDispatch` ilə komponentlərə bağlanır. Boilerplate əhəmiyyətli dərəcədə azalıb.

## Əsas terminlər

- Context API
- createContext
- Provider
- useContext
- SWR
- Stale-While-Revalidate
- Suspense
- Redux Toolkit
- createSlice
- configureStore
- useSelector
- useDispatch
- Skeleton
- Caching

## Praktik nəticə

Pokedex nümunəsini SWR + Suspense ilə qur: 151 Pokemon paralel yüklə, skeleton göstər, sonra grid göstər.

## Mənbə

Pages: 235-258
# 9. React Router

**Səhifələr:** 192-216

## Bu fəsil nədən bəhs edir?

React Router SPA-larda (single-page application) URL əsaslı naviqasiya üçün istifadə olunan de-facto kitabxanadır. Bu fəsil v6.4-ə qədər və v6.4-dəki yeni data API-ni əhatə edir.

## Əsas fikirlər

### 1. React Router versiyaları

`react-router` — ümumi paket; `react-router-dom` — web üçün; `react-router-native` — React Native üçün. Browser tətbiqlərində adətən `react-router-dom` istifadə olunur.

### 2. Marşrut qurma

`createBrowserRouter([{ path, element, children }])` ilə marşrutlar təyin olunur. `<RouterProvider router={router} />` ilə tətbiqə inject olunur. `<Routes>` komponenti içində `<Route path=... element=... />` istifadə olunur.

### 3. URL parametrləri

`<Route path='/users/:id'>` yazanda `useParams()` hook-u ilə URL-dən oxud. Məsələn, `/users/42` → `{ id: '42' }`. Parametr dinamik olduğu üçün `useEffect` ilə yeni ID-yə görə data fetch etmək lazımdır.

### 4. React Router v6.4 yenilikləri

v6.4 data API təqdim etdi: `loader`, `action`, `defer`, `errorElement`. `loader` marşrut yüklənməmişdən əvvəl data fetch edir — bu, waterfall problemi həll edir. `<RouterProvider>` ilə birlikdə işləyir.

### 5. Loaders

`loader: async ({ params }) => fetch(`/api/users/${params.id}`)` yazanda, komponent render olunmazdan əvvəl data hazır olur. Komponentdə `useLoaderData()` ilə data oxunur. Bu pattern Suspense ilə birlikdə yaxşı işləyir.

## Əsas terminlər

- React Router
- SPA
- Routing
- createBrowserRouter
- Routes
- Route
- params
- useParams
- Loader
- useLoaderData
- Suspense
- data API

## Praktik nəticə

3 səhifəli (Home, About, Contact) kiçik bir SPA qur React Router ilə. Sonra bir səhifəyə loader əlavə et və mock API-dən istifadəçi məlumatı yüklə.

## Mənbə

Pages: 192-216
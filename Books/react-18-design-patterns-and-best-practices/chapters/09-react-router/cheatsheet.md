# Cheat Sheet — React Router

## React Router

### BrowserRouter ilə marşrut

**Nə edir:** `createBrowserRouter` ilə marşrutlar, `RouterProvider` ilə tətbiqə inject.

**Kod:**

```tsx
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
const router = createBrowserRouter([
  { path: '/', element: <Home /> },
  { path: '/about', element: <About /> },
]);
<RouterProvider router={router} />
```

**Mənbə:** Chapter 9, page 192-216

### Link ilə naviqasiya

**Nə edir:** `a` teq əvəzinə — SPA naviqasiya üçün, brauzerin refresh etmir.

**Kod:**

```tsx
import { Link } from 'react-router-dom';
<Link to='/about'>About</Link>
```

**Mənbə:** Chapter 9, page 192-216

### URL params

**Nə edir:** Dinamik URL hissələrini oxumaq üçün `useParams`.

**Kod:**

```tsx
// Route
{ path: 'users/:id', element: <UserDetail /> }
// Komponent
import { useParams } from 'react-router-dom';
const { id } = useParams();
```

**Mənbə:** Chapter 9, page 192-216

### Loader (v6.4)

**Nə edir:** Marşrut render olunmazdan əvvəl data fetch edir.

**Kod:**

```tsx
{
  path: 'users/:id',
  element: <UserDetail />,
  loader: async ({ params }) =>
    fetch(`/api/users/${params.id}`).then(r => r.json()),
}
// Komponentdə
const user = useLoaderData() as User;
```

**Mənbə:** Chapter 9, page 192-216

# Cheat Sheet — React 18 New Features

## React 18

### createRoot

**Nə edir:** Köhnə `ReactDOM.render` əvəzinə React 18-in yeni entry point.

**Kod:**

```ts
import { createRoot } from 'react-dom/client';
const root = createRoot(document.getElementById('root')!);
root.render(<App />);
```

**Mənbə:** Chapter 10, page 217-234

### useTransition

**Nə edir:** Aşağı prioritetli yeniləmələri işarələyir. `isPending` loading indikatoru üçün.

**Kod:**

```tsx
const [isPending, startTransition] = useTransition();
startTransition(() => {
  setQuery(input);
});
```

**Mənbə:** Chapter 10, page 217-234

### useDeferredValue

**Nə edir:** Axtarış input kimi ssenarilərdə siyahı yeniləməsini gecikdirmək üçün.

**Kod:**

```tsx
const query = useDeferredValue(searchInput);
```

**Mənbə:** Chapter 10, page 217-234

### useId

**Nə edir:** SSR-uyğun unikal ID. Form label-id əlaqəsi üçün ideal.

**Kod:**

```tsx
const id = useId();
<label htmlFor={id}>Name</label>
<input id={id} />
```

**Mənbə:** Chapter 10, page 217-234

### useInsertionEffect

**Nə edir:** CSS-in-JS kitabxanaları üçün, useEffect-dən daha erkən çağırılır.

**Kod:**

```tsx
useInsertionEffect(() => {
  const style = document.createElement('style');
  style.textContent = rules;
  document.head.appendChild(style);
}, [rules]);
```

**Mənbə:** Chapter 10, page 217-234

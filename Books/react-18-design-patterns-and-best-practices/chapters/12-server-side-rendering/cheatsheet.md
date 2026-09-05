# Cheat Sheet — Server-Side Rendering

## SSR

### renderToString (Express)

**Nə edir:** Server-də React HTML-ə çevrilir.

**Kod:**

```ts
import { renderToString } from 'react-dom/server';
app.get('/', (req, res) => {
  const html = renderToString(<App />);
  res.send(`<!DOCTYPE html><div id='root'>${html}</div>`);
});
```

**Mənbə:** Chapter 12, page 259-280

### hydrateRoot

**Nə edir:** Server HTML-i client-də canlandırır (hydration).

**Kod:**

```ts
import { hydrateRoot } from 'react-dom/client';
hydrateRoot(document.getElementById('root')!, <App />);
```

**Mənbə:** Chapter 12, page 259-280

### getServerSideProps

**Nə edir:** Next.js-də hər sorğuda server-də data fetch.

**Kod:**

```tsx
export const getServerSideProps: GetServerSideProps = async () => {
  const res = await fetch('https://api.example.com/data');
  const data = await res.json();
  return { props: { data } };
};
```

**Mənbə:** Chapter 12, page 259-280

### getStaticProps

**Nə edir:** Build zamanı HTML yaradır, ISR (Incremental Static Regeneration).

**Kod:**

```tsx
export const getStaticProps = async () => {
  const data = await getData();
  return { props: { data }, revalidate: 60 };
};
```

**Mənbə:** Chapter 12, page 259-280

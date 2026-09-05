# Cheat Sheet — React 18 Design Patterns and Best Practices

## JSX Sintaksis

### Sadə JSX
```jsx
const element = <h1>Hello, world!</h1>;
```

### Şərti render
```jsx
{isLoggedIn ? <Dashboard /> : <Login />}
{showWarning && <Warning />}
```

### Siyahı render
```jsx
{items.map((item) => (
  <Item key={item.id} {...item} />
))}
```

### Çox xətt
```jsx
return (
  <div>
    <h1>Title</h1>
    <p>Paragraph</p>
  </div>
);
```

## React 18 Yeni API-ləri

### createRoot
```jsx
import { createRoot } from 'react-dom/client';

const container = document.getElementById('root');
const root = createRoot(container);
root.render(<App />);
```

### hydrateRoot
```jsx
import { hydrateRoot } from 'react-dom/client';

hydrateRoot(container, <App />);
```

### renderToPipeableStream (SSR)
```jsx
import { renderToPipeableStream } from 'react-dom/server';

const { pipe, abort } = renderToPipeableStream(<App />, {
  onShellReady() {
    response.setHeader('Content-Type', 'text/html');
    pipe(response);
  },
});
```

## Yeni Hooks

### useId — Unikal ID generasiyası
```jsx
const id = useId();
return (
  <>
    <label htmlFor={id}>Name</label>
    <input id={id} type="text" />
  </>
);
```

### useTransition — Prioritetli yeniləmələr
```jsx
const [isPending, startTransition] = useTransition();

const handleClick = () => {
  startTransition(() => {
    setTab('profile');
  });
};
```

### useDeferredValue — Gecikdirilmiş dəyər
```jsx
const deferredQuery = useDeferredValue(query);
```

### useInsertionEffect — CSS injection
```jsx
useInsertionEffect(() => {
  const style = document.createElement('style');
  style.textContent = rules;
  document.head.appendChild(style);
}, [rules]);
```

## Əsas Hooks

### useState
```jsx
const [count, setCount] = useState(0);
setCount(count + 1);
setCount((prev) => prev + 1);
```

### useEffect
```jsx
useEffect(() => {
  // side effect
  const sub = api.subscribe(id, setData);

  return () => sub.unsubscribe();
}, [id]);
```

### useMemo — Dəyəri memoizasiya
```jsx
const sorted = useMemo(() => items.sort(compareFn), [items]);
```

### useCallback — Funksiyanı memoizasiya
```jsx
const handleClick = useCallback((id) => {
  setItems((prev) => prev.filter((x) => x.id !== id));
}, []);
```

### useReducer
```jsx
function reducer(state, action) {
  switch (action.type) {
    case 'increment': return { count: state.count + 1 };
    case 'decrement': return { count: state.count - 1 };
    default: return state;
  }
}

const [state, dispatch] = useReducer(reducer, { count: 0 });
dispatch({ type: 'increment' });
```

## TypeScript ilə React

### Komponent tipi
```tsx
import { FC, PropsWithChildren } from 'react';

type Props = PropsWithChildren<{ title: string; count?: number }>;

const Card: FC<Props> = ({ title, count = 0, children }) => (
  <div>
    <h2>{title}</h2>
    <span>{count}</span>
    {children}
  </div>
);
```

### useState ilə tip
```tsx
const [user, setUser] = useState<User | null>(null);
const [items, setItems] = useState<Item[]>([]);
```

### useRef ilə tip
```tsx
const inputRef = useRef<HTMLInputElement>(null);

inputRef.current?.focus();
```

### forwardRef ilə tip
```tsx
const MyInput = forwardRef<HTMLInputElement, Props>((props, ref) => (
  <input ref={ref} {...props} />
));
```

## Komponent Patternləri

### Container / Presentational
```jsx
// Container — məntiq
const UserListContainer: FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  // ... data fetching logic
  return <UserList users={users} />;
};

// Presentational — UI
const UserList: FC<{ users: User[] }> = ({ users }) => (
  <ul>{users.map((u) => <li key={u.id}>{u.name}</li>)}</ul>
);
```

### HOC (Higher-Order Component)
```tsx
function withAuth<P>(Component: FC<P>) {
  return (props: P) => {
    const { user } = useAuth();
    if (!user) return <Redirect to="/login" />;
    return <Component {...props} />;
  };
}
```

### Function as Child
```jsx
<Mouse position={(x, y) => <Cursor x={x} y={y} />}>
  {({ x, y }) => <Cursor x={x} y={y} />}
</Mouse>
```

## Form İdarəetməsi

### Controlled
```jsx
const [name, setName] = useState('');
return <input value={name} onChange={(e) => setName(e.target.value)} />;
```

### Uncontrolled (ref ilə)
```jsx
const ref = useRef<HTMLInputElement>(null);
const handleSubmit = () => {
  console.log(ref.current?.value);
};
return (
  <form onSubmit={handleSubmit}>
    <input ref={ref} defaultValue="" />
  </form>
);
```

## Stil üsulları

### Inline style
```jsx
<div style={{ color: 'red', fontSize: '1.2rem' }}>Text</div>
```

### CSS Modules
```css
/* Button.module.css */
.primary { background: blue; }
```
```jsx
import styles from './Button.module.css';
<button className={styles.primary}>Click</button>
```

### styled-components
```tsx
import styled from 'styled-components';

const Button = styled.button`
  background: ${(p) => (p.primary ? 'blue' : 'white')};
  padding: 0.5rem 1rem;
`;

<Button primary>Click</Button>
```

## React Router v6.4

### Routing
```jsx
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

const router = createBrowserRouter([
  {
    path: '/',
    element: <Root />,
    loader: rootLoader,
    children: [
      { index: true, element: <Home /> },
      { path: 'contact/:id', element: <Contact /> },
    ],
  },
]);

<RouterProvider router={router} />
```

### Loaders
```jsx
export async function rootLoader() {
  const contacts = await getContacts();
  return { contacts };
}
```

## Data Management

### Context API
```tsx
const ThemeContext = createContext<Theme>('light');

function App() {
  const [theme, setTheme] = useState<Theme>('light');
  return (
    <ThemeContext.Provider value={{ theme, setTheme }}>
      <Page />
    </ThemeContext.Provider>
  );
}

function Page() {
  const { theme } = useContext(ThemeContext);
  return <div className={theme}>Content</div>;
}
```

### SWR (data fetching)
```jsx
import useSWR from 'swr';

const fetcher = (url) => fetch(url).then((r) => r.json());

function Pokemon({ name }) {
  const { data, error } = useSWR(`/api/pokemon/${name}`, fetcher);
  if (error) return <div>failed</div>;
  if (!data) return <div>loading...</div>;
  return <div>{data.name}</div>;
}
```

### Redux Toolkit Slice
```tsx
import { createSlice, configureStore } from '@reduxjs/toolkit';

const counterSlice = createSlice({
  name: 'counter',
  initialState: { value: 0 },
  reducers: {
    increment: (state) => { state.value += 1 },
    decrement: (state) => { state.value -= 1 },
  },
});

const store = configureStore({ reducer: counterSlice.reducer });
```

## GraphQL (Apollo)

### Query
```tsx
import { gql, useQuery } from '@apollo/client';

const GET_USERS = gql`
  query GetUsers {
    users { id name email }
  }
`;

function Users() {
  const { loading, error, data } = useQuery(GET_USERS);
  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error :(</p>;
  return data.users.map((u) => <p key={u.id}>{u.name}</p>);
}
```

### Mutation
```tsx
const LOGIN = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) { token user { id name } }
  }
`;

const [login, { data, loading, error }] = useMutation(LOGIN);
```

## Testing — Vitest

### Sadə test
```tsx
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import Counter from './Counter';

describe('Counter', () => {
  it('renders initial count', () => {
    render(<Counter initial={5} />);
    expect(screen.getByText('5')).toBeInTheDocument();
  });

  it('increments on click', async () => {
    render(<Counter initial={0} />);
    await userEvent.click(screen.getByRole('button'));
    expect(screen.getByText('1')).toBeInTheDocument();
  });
});
```

## Server-Side Rendering

### Next.js ilə səhifə
```tsx
import { GetServerSideProps } from 'next';

export const getServerSideProps: GetServerSideProps = async () => {
  const res = await fetch('https://api.example.com/data');
  const data = await res.json();
  return { props: { data } };
};

export default function Page({ data }) {
  return <div>{data.title}</div>;
}
```

## Anti-patterns (qarşısını alın)

### ❌ Index key
```jsx
{items.map((item, index) => <Item key={index} {...item} />)}
```
✅ Doğrusu:
```jsx
{items.map((item) => <Item key={item.id} {...item} />)}
```

### ❌ State-i prop ilə initializasiya
```jsx
function Counter({ initial }) {
  const [count, setCount] = useState(initial); // ❌ prop dəyişsə state yenilənməz
}
```
✅ Doğrusu:
```jsx
function Counter({ initial }) {
  const [count, setCount] = useState(initial);
  // və ya initialKey istifadə edin
}
```

### ❌ Qeyri-standart prop-ları DOM-a spread
```jsx
<div {...props}> // ❌ React tanımayan atributlar DOM-a DOM warning yaradır
```
✅ Doğrusu:
```jsx
const { known, ...rest } = props;
<div {...rest}> // yalnız DOM atributları
```

## Performans Optimallaşdırması

### React.memo ilə komponent memoizasiyası
```tsx
const ExpensiveList = React.memo(({ items }) => (
  <ul>{items.map((i) => <Item key={i.id} {...i} />)}</ul>
));
```

### Reconciliation üçün düzgün key
```jsx
// Çox elementli siyahılarda həmişə unikal, sabit ID istifadə edin
{users.map((u) => <UserCard key={u.id} user={u} />)}
```

### Babel Plugin: lodash → individual imports
```js
// .babelrc
{
  "plugins": [
    ["lodash", { "id": "lodash" }],
    ["import", { "libraryName": "antd" }]
  ]
}
```

## Deployment

### Nginx konfiqurasiyası (reverse proxy)
```nginx
server {
  listen 80;
  server_name example.com;

  location / {
    proxy_pass http://localhost:3000;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection 'upgrade';
    proxy_set_header Host $host;
    proxy_cache_bypass $http_upgrade;
  }
}
```

### PM2 ilə Node.js
```bash
pm2 start npm --name "app" -- start
pm2 save
pm2 startup
```

### CircleCI konfiqurasiyası
```yaml
version: 2.1
jobs:
  build:
    docker:
      - image: cimg/node:18.0
    steps:
      - checkout
      - run: npm ci
      - run: npm run lint
      - run: npm test
      - run: npm run build
```

## MonoRepo (NPM Workspaces)

### package.json (root)
```json
{
  "name": "my-monorepo",
  "workspaces": ["packages/*", "apps/*"]
}
```

### Workspace əmrləti
```bash
npm install lodash -w @my-org/utils
npm run build -ws --if-present
npm run test -ws --if-present
```

## Faydalı Qaydalar

- `key` props üçün unikal, sabit ID istifadə edin (index yox)
- Hooks yalnız React funksiyalarının top-level-ində çağırılmalıdır
- State və props həmişə immutable qalsın
- Komponent adları böyük hərflə başlayır (`<MyComponent />`)
- Mümkün yerlərdə functional komponent istifadə edin (class komponent köhnəlib)
- React DevTools həmişə açıq olsun development zamanı
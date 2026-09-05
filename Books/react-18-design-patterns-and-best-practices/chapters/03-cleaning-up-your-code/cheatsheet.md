# Cheat Sheet — Cleaning Up Your Code

## JSX, Babel, ESLint, FP

### JSX element

**Nə edir:** HTML-ə oxşar sintaksis. `className` (class yox), camelCase attributes.

**Kod:**

```jsx
const el = <div className='box'>Hello</div>;
```

**Mənbə:** Chapter 3, page 58-80

### Şərti render

**Nə edir:** `&&` qısa yoludur — yalnız `true` olduqda render edir.

**Kod:**

```jsx
{isLoggedIn ? <Dashboard /> : <Login />}
{show && <Modal />}
```

**Mənbə:** Chapter 3, page 58-80

### Siyahı

**Nə edir:** `key` prop-u unikal olmalıdır (index yox!).

**Kod:**

```jsx
{users.map(u => (
  <Card key={u.id} {...u} />
))}
```

**Mənbə:** Chapter 3, page 58-80

### Babel preset

**Nə edir:** Babel müasir JSX/TS-i köhnə JavaScript-ə çevirir.

**Kod:**

```json
{
  "presets": ["@babel/preset-env", "@babel/preset-react", "@babel/preset-typescript"]
}
```

**Mənbə:** Chapter 3, page 58-80

### ESLint flat config

**Nə edir:** ESLint qaydaları təyin edir, səhvləri compile-time-da tutur.

**Kod:**

```js
export default [
  { rules: { 'react-hooks/rules-of-hooks': 'error' } }
];
```

**Mənbə:** Chapter 3, page 58-80

### Pure function

**Nə edir:** Eyni giriş həmişə eyni çıxışı verir, yan təsirsiz.

**Kod:**

```ts
function add(a: number, b: number): number {
  return a + b;
}
```

**Mənbə:** Chapter 3, page 58-80

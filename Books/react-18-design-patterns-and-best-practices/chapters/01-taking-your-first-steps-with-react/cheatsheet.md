# Cheat Sheet — Taking Your First Steps with React

## İlk React addımları

### İlk komponent

**Nə edir:** Funksional komponent, JSX ilə return edir. `<App />` kimi istifadə olunur.

**Kod:**

```tsx
function App() {
  return <h1>Hello, React!</h1>;
}
```

**Mənbə:** Chapter 1, page 30-43

### useState ilə state

**Nə edir:** `useState` hook-u komponentə state əlavə edir. `count` dəyər, `setCount` yeniləyici funksiyadır.

**Kod:**

```tsx
const [count, setCount] = useState(0);
<button onClick={() => setCount(count + 1)}>
  Count: {count}
</button>
```

**Mənbə:** Chapter 1, page 30-43

### Vite ilə yeni layihə

**Nə edir:** Vite sürətli development server təmin edir, HMR dəstəkləyir.

**Kod:**

```shell
npm create vite@latest my-app -- --template react-ts
```

**Mənbə:** Chapter 1, page 30-43

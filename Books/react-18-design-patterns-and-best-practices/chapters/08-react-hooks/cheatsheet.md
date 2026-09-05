# Cheat Sheet — React Hooks

## React Hooks

### useState

**Nə edir:** `setValue`-a funksiya ötürmək async yeniləmələrdə köhnə state-ə görə işləyir.

**Kod:**

```tsx
const [value, setValue] = useState(0);
setValue(prev => prev + 1);
```

**Mənbə:** Chapter 8, page 153-191

### useEffect

**Nə edir:** Yan təsir həyata keçirir. Cleanup return edir. `deps` massivi vacibdir.

**Kod:**

```tsx
useEffect(() => {
  const sub = api.subscribe(id, setData);
  return () => sub.unsubscribe();
}, [id]);
```

**Mənbə:** Chapter 8, page 153-191

### useMemo

**Nə edir:** Ağır hesablamaları memoizasiya edir, deps dəyişəndə yenidən hesablayır.

**Kod:**

```tsx
const sorted = useMemo(
  () => items.sort(comparator),
  [items]
);
```

**Mənbə:** Chapter 8, page 153-191

### useCallback

**Nə edir:** Funksiya referensini memoizasiya edir, uşaq komponentə prop kimi ötürüləndə faydalıdır.

**Kod:**

```tsx
const onClick = useCallback((id: string) => {
  setItems(prev => prev.filter(x => x.id !== id));
}, []);
```

**Mənbə:** Chapter 8, page 153-191

### useReducer

**Nə edir:** Complex state-lər üçün. Reducer pure funksiyadır: `(state, action) => newState`.

**Kod:**

```tsx
const [state, dispatch] = useReducer(reducer, { count: 0 });
dispatch({ type: 'INC' });
```

**Mənbə:** Chapter 8, page 153-191

### Custom Hook

**Nə edir:** `use` ilə başlayır, digər hook-ları çağıra bilər, stateful məntiqi paylaşır.

**Kod:**

```tsx
function useWindowWidth() {
  const [w, setW] = useState(window.innerWidth);
  useEffect(() => {
    const onResize = () => setW(window.innerWidth);
    window.addEventListener('resize', onResize);
    return () => window.removeEventListener('resize', onResize);
  }, []);
  return w;
}
```

**Mənbə:** Chapter 8, page 153-191

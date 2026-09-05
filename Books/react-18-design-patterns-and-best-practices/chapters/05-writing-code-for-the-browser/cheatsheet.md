# Cheat Sheet — Writing Code for the Browser

## Forms, Events, Refs

### Controlled input

**Nə edir:** Dəyər React state-də, hər dəyişiklikdə state yenilənir.

**Kod:**

```jsx
const [name, setName] = useState('');
<input value={name} onChange={e => setName(e.target.value)} />
```

**Mənbə:** Chapter 5, page 94-114

### Uncontrolled input (ref ilə)

**Nə edir:** Dəyər DOM-da qalır, lazım olanda ref ilə oxunur.

**Kod:**

```jsx
const ref = useRef<HTMLInputElement>(null);
const onSubmit = () => console.log(ref.current?.value);
<input ref={ref} defaultValue='' />
```

**Mənbə:** Chapter 5, page 94-114

### Event handler

**Nə edir:** Event adları camelCase. SyntheticEvent React tərəfindən idarə olunur.

**Kod:**

```jsx
<button onClick={handleClick}>Click</button>
```

**Mənbə:** Chapter 5, page 94-114

### useRef

**Nə edir:** DOM elementinə istinad. State kimi render-ə səbəb olmur.

**Kod:**

```tsx
const ref = useRef<HTMLDivElement>(null);
useEffect(() => { ref.current?.focus(); }, []);
```

**Mənbə:** Chapter 5, page 94-114

### forwardRef

**Nə edir:** Ref-i uşaq komponentdən dərin DOM elementə yönləndirir.

**Kod:**

```tsx
const Input = forwardRef<HTMLInputElement, Props>((props, ref) => (
  <input ref={ref} {...props} />
));
```

**Mənbə:** Chapter 5, page 94-114

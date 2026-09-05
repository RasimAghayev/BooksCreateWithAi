# Cheat Sheet — Managing Data

## Data Management

### Context API

**Nə edir:** Props drilling həll edir, amma çox consumer olduqda performans problemi yaradır.

**Kod:**

```tsx
const ThemeContext = createContext('light');
function App() {
  const [theme, setTheme] = useState('light');
  return (
    <ThemeContext.Provider value={{ theme, setTheme }}>
      <Page />
    </ThemeContext.Provider>
  );
}
// İstehlakçı
const { theme } = useContext(ThemeContext);
```

**Mənbə:** Chapter 11, page 235-258

### SWR

**Nə edir:** Stale-while-revalidate caching strategiyası.

**Kod:**

```tsx
import useSWR from 'swr';
const fetcher = (url: string) => fetch(url).then(r => r.json());
function Profile() {
  const { data, error } = useSWR('/api/user', fetcher);
  if (error) return <div>failed</div;
  if (!data) return <div>loading...</div>;
  return <div>{data.name}</div>;
}
```

**Mənbə:** Chapter 11, page 235-258

### Redux Toolkit Slice

**Nə edir:** Modern Redux, az boilerplate, immer ilə immutable yeniləmələr.

**Kod:**

```tsx
const counterSlice = createSlice({
  name: 'counter',
  initialState: { value: 0 },
  reducers: {
    increment: state => { state.value += 1 },
  },
});
export const { increment } = counterSlice.actions;
export default counterSlice.reducer;
```

**Mənbə:** Chapter 11, page 235-258

### useSelector / useDispatch

**Nə edir:** Komponentdən Redux state-ə giriş və dispatch.

**Kod:**

```tsx
const count = useSelector((s: RootState) => s.counter.value);
const dispatch = useDispatch();
<button onClick={() => dispatch(increment())}>+</button>
```

**Mənbə:** Chapter 11, page 235-258

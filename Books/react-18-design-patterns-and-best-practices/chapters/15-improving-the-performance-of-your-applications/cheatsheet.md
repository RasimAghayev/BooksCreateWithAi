# Cheat Sheet — Improving the Performance of Your Applications

## Performance

### React.memo

**Nə edir:** Props dəyişmədikcə yenidən render etmir.

**Kod:**

```tsx
const Card = React.memo(({ user }: { user: User }) => (
  <div>{user.name}</div>
));
```

**Mənbə:** Chapter 15, page 413-420

### Lazy loading

**Nə edir:** Komponent tələb olunanda yüklənir.

**Kod:**

```tsx
const Heavy = React.lazy(() => import('./Heavy'));
<Suspense fallback={<Spinner />}>
  <Heavy />
</Suspense>
```

**Mənbə:** Chapter 15, page 413-420

### useImmer (Immer ilə)

**Nə edir:** Immer ilə asanlıqla immutable state yeniləmələri.

**Kod:**

```tsx
import { useImmer } from 'use-immer';
const [user, updateUser] = useImmer({ name: '', age: 0 });
updateUser(d => { d.age = 30; });
```

**Mənbə:** Chapter 15, page 413-420

### react-window (virtualization)

**Nə edir:** Böyük siyahıları pəncərələmək — DOM-da yalnız görünən elementlər.

**Kod:**

```tsx
import { FixedSizeList } from 'react-window';
<FixedSizeList height={600} itemCount={10000} itemSize={35}>
  {({ index, style }) => <div style={style}>Row {index}</div>}
</FixedSizeList>
```

**Mənbə:** Chapter 15, page 413-420

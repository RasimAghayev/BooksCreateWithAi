# Cheat Sheet — Exploring Popular Composition Patterns

## Composition Patterns

### Container / Presentational

**Nə edir:** Məntiq və UI-ı ayırır, presentational komponentlər təkrar istifadə edilə bilər.

**Kod:**

```tsx
// Container
const UserListContainer = () => {
  const [users, setUsers] = useState<User[]>([]);
  useEffect(() => { fetchUsers().then(setUsers); }, []);
  return <UserList users={users} />;
};
// Presentational
const UserList: FC<{ users: User[] }> = ({ users }) => (
  <ul>{users.map(u => <li key={u.id}>{u.name}</li>)}</ul>
);
```

**Mənbə:** Chapter 4, page 81-93

### HOC nümunəsi

**Nə edir:** HOC komponentə cross-cutting concern əlavə edir. Hooks ilə daha az istifadə olunur.

**Kod:**

```tsx
function withLoading<P>(Component: FC<P>) {
  return (props: P) => {
    const { loading } = useLoading();
    if (loading) return <Spinner />;
    return <Component {...props} />;
  };
}
const LoadedProfile = withLoading(Profile);
```

**Mənbə:** Chapter 4, page 81-93

### FunctionAsChild

**Nə edir:** `children` prop kimi funksiya göndərilir — render prop pattern.

**Kod:**

```jsx
<ThemeConsumer>
  {({ theme }) => (
    <div className={theme}>Content</div>
  )}
</ThemeConsumer>
```

**Mənbə:** Chapter 4, page 81-93

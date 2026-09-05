# Cheat Sheet — Anti-Patterns to Be Avoided

## Anti-patterns

### ❌ Index key

**Nə edir:** `key` üçün unikal, sabit ID istifadə edin. Index reconciliation pozur.

**Kod:**

```jsx
{items.map((item, i) => <Item key={i} {...item} />)}
```

**Mənbə:** Chapter 7, page 142-152

### ✅ Unikal key

**Nə edir:** React elementi düzgün identifikasiya edir.

**Kod:**

```jsx
{items.map(item => <Item key={item.id} {...item} />)}
```

**Mənbə:** Chapter 7, page 142-152

### ❌ State-i prop ilə

**Nə edir:** Prop dəyişsə state yenilənmir. State-i valideynə qaldırın.

**Kod:**

```jsx
function Counter({ initial }) {
  const [count, setCount] = useState(initial);
}
```

**Mənbə:** Chapter 7, page 142-152

### ✅ DOM spread zamanı filtr

**Nə edir:** Yalnız DOM-a aid atributları ötürün.

**Kod:**

```jsx
const { known, ...rest } = props;
<div {...rest}>...</div>
```

**Mənbə:** Chapter 7, page 142-152

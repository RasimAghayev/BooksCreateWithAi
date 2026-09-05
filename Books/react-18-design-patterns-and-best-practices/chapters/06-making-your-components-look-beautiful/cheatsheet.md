# Cheat Sheet — Making Your Components Look Beautiful

## Stil yanaşmaları

### Inline style

**Nə edir:** Lokal və dinamik, lakin media queries və pseudo-classes dəstəkləmir.

**Kod:**

```jsx
<div style={{ color: 'red', fontSize: 16 }}>Text</div>
```

**Mənbə:** Chapter 6, page 115-141

### CSS Modules

**Nə edir:** Lokal scoped CSS, Webpack 5 default dəstəkləyir.

**Kod:**

```jsx
// Card.module.css
.title { font-size: 1.5rem; }

// komponent:
import s from './Card.module.css';
<h2 className={s.title}>Title</h2>
```

**Mənbə:** Chapter 6, page 115-141

### styled-components

**Nə edir:** Runtime CSS-in-JS, dinamik prop-larla stil dəyişir.

**Kod:**

```tsx
const Button = styled.button<{ primary?: boolean }>`
  background: ${p => p.primary ? 'blue' : 'white'};
  padding: 0.5rem 1rem;
`;
```

**Mənbə:** Chapter 6, page 115-141

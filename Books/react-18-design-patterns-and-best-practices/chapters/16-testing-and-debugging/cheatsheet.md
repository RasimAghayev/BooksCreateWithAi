# Cheat Sheet — Testing and Debugging

## Testing

### Vitest əsas test

**Nə edir:** Vitest API Jest-ə uyğundur, daha sürətlidir.

**Kod:**

```ts
import { describe, it, expect } from 'vitest';
import { add } from './math';
describe('add', () => {
  it('sums two numbers', () => {
    expect(add(2, 3)).toBe(5);
  });
});
```

**Mənbə:** Chapter 16, page 421-466

### React Testing Library — render

**Nə edir:** Komponenti render et, screen ilə DOM-u sorğula.

**Kod:**

```tsx
import { render, screen } from '@testing-library/react';
import Counter from './Counter';
it('shows initial count', () => {
  render(<Counter initial={5} />);
  expect(screen.getByText('5')).toBeInTheDocument();
});
```

**Mənbə:** Chapter 16, page 421-466

### userEvent ilə klik

**Nə edir:** Real istifadəçi davranışını simulyasiya et.

**Kod:**

```tsx
import userEvent from '@testing-library/user-event';
it('increments on click', async () => {
  render(<Counter initial={0} />);
  await userEvent.click(screen.getByRole('button'));
  expect(screen.getByText('1')).toBeInTheDocument();
});
```

**Mənbə:** Chapter 16, page 421-466

### Mock funksiya

**Nə edir:** `vi.fn()` mock funksiya yaradır, çağırışları yoxlamaq olur.

**Kod:**

```ts
const onSubmit = vi.fn();
render(<Form onSubmit={onSubmit} />);
await userEvent.click(screen.getByText('Submit'));
expect(onSubmit).toHaveBeenCalledTimes(1);
```

**Mənbə:** Chapter 16, page 421-466

### Snapshot test

**Nə edir:** Render nəticəsini snapshot ilə müqayisə et.

**Kod:**

```tsx
it('matches snapshot', () => {
  const { container } = render(<Card title='Hi' />);
  expect(container).toMatchSnapshot();
});
```

**Mənbə:** Chapter 16, page 421-466

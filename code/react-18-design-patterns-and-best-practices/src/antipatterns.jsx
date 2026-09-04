// Anti-patterns to avoid: state-from-props, index keys, spreading DOM props
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 7 (pp. 151-157)

import { FC, useState } from 'react';

// ❌ ANTI-PATTERN: initialize state from props → duplicated source of truth + stale when prop changes
export const CounterBad = ({ count }) => {
  // useState runs once; if `count` prop later changes, state stays stale (1→2 in props, state still 1)
  const [state, setState] = useState({ count });
  const handleClick = () => setState({ count: state.count + 1 });
  return (
    <div>
      {state.count}
      <button onClick={handleClick}>+</button>
    </div>
  );
};

// ✅ FIX: name the prop explicitly (initialCount) to signal it is one-time
type Props = { initialCount: number };
export const CounterGood: FC<Props> = ({ initialCount }) => {
  const [count, setCount] = useState(initialCount);
  const handleClick = () => setCount(count + 1);
  return (
    <div>
      {count}
      <button onClick={handleClick}>+</button>
    </div>
  );
};

// ❌ ANTI-PATTERN: using .map index as key → React reuses nodes by position on insert
const ListBad = () => {
  const [items, setItems] = useState(['foo', 'bar']);
  const handleClick = () => setItems(['baz', ...items]);
  return (
    <div>
      <ul>
        {items.map((item, index) => (
          <li key={index}>
            {item} <input type="text" />
          </li>
        ))}
      </ul>
      <button onClick={handleClick}>+</button>
    </div>
  );
};

// ✅ FIX: key must be unique + stable. Combine content+index when content may repeat.
export const ListGood = () => {
  const [items, setItems] = useState(['foo', 'bar']);
  const handleClick = () => setItems(['baz', ...items]);
  return (
    <div>
      <ul>
        {items.map((item, index) => (
          <li key={`${item}-${index}`}>
            {item} <input type="text" />
          </li>
        ))}
      </ul>
      <button onClick={handleClick}>+</button>
    </div>
  );
};

// ❌ ANTI-PATTERN: spreading arbitrary props onto a DOM element → passes unknown attrs (`foo`)
const SpreadBad = props => <div {...props} />;
// Usage: <SpreadBad className="baz" foo="bar" /> → "Unknown prop `foo` on <div> tag" warning

// ✅ FIX: explicitly separate DOM-valid props
export const SpreadGood = ({ className, ...domProps }) => (
  <div className={className} {...domProps} />
);

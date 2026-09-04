// JSX expressions: attributes, children, style, fragments, conditionals
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 3

const styles = {
  backgroundColor: 'red',
};

// Props (attributes)
const img = _jsx('img', {
  src: 'https://www.ranchosanpancho.com/images/logo.png',
  alt: 'Cabañas San Pancho',
});

// Children + nested elements
const linkTree = _jsx(
  'div',
  null,
  _jsx('a', { href: 'https://ranchosanpancho.com' }, 'Click me!')
);

// HTML vs JSX attribute differences
const labelEl = _jsx('label', {
  className: 'awesome-label',
  htmlFor: 'name',
});

// Inline style as object (camelCase)
const styledDiv = _jsx('div', { style: styles });

// React.Fragment for wrapping without extra DOM node
import { Fragment } from 'react';
function App() {
  return (
    <Fragment>
      <h1>An h1 heading</h1>
      Some text here.
      <h2>An h2 heading</h2>
    </Fragment>
  );
}

// Shorthand fragment syntax
function AppShort() {
  return (
    <>
      <ComponentA />
      <ComponentB />
      <ComponentC />
    </>
  );
}

// Explicit spaces in JSX (JSX collapses whitespace)
function NameWithSpace() {
  return (
    <div>
      <span>My</span>
      {' '}
      name is
      {' '}
      <span>Carlos</span>
    </div>
  );
}

// Boolean attributes (omitted = true; must be explicit for false)
const disabledTrue = _jsx('button', { disabled: true });
const disabledFalse = _jsx('button', { disabled: false });

// Spread attributes
const attrs = {
  id: 'myId',
  className: 'myClass',
};
const spreadDiv = _jsx('div', attrs);

// Template literals
const name = 'Carlos';
const age = 35;
const message = `Hello, my name is ${name} and I am ${age} years old.`;

// Multiline JSX with parentheses (avoid ASI issues)
function Render() {
  return (
    <div>
      <Header />
      <div>
        <Main content={...} />
      </div>
    </div>
  );
}

// Multi-property element (aligned attributes)
function BigButton({ onSomething }) {
  return (
    <button
      foo="bar"
      veryLongPropertyName="baz"
      onSomething={onSomething}
    />
  );
}

// Conditional rendering: && operator
function UserSession({ isLoggedIn }) {
  return (
    <div>
      {isLoggedIn && <LogoutButton />}
    </div>
  );
}

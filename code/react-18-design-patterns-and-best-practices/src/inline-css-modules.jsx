// Inline styles in React + dynamic styles; CSS modules scoped import
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 6 (pp. 126-139)

import { useState, ChangeEvent } from 'react';

// Inline styles: camelCased CSS keys, string values in quotes (numbers = px)
const buttonStyle = {
  color: 'palevioletred',
  backgroundColor: 'papayawhip',
};

export const InlineButton = () => (
  <button style={buttonStyle}>Click me!</button>
);

// Vendor prefixes: capitalize (except 'ms')
const webkitTransition = { WebkitTransition: 'all 0.3s ease' };

// Numbers → px by default
const fixedHeight = { height: 100 }; // 100px

// Dynamic inline style: font-size changes with input value
export const FontSize = () => {
  const [value, setValue] = useState(16);

  const handleChange = (e) => setValue(Number(e.target.value));

  return (
    <input
      type="number"
      value={value}
      onChange={handleChange}
      style={{ fontSize: value }}
    />
  );
};

// CSS modules: class names are locally scoped to a hash
import styles from './index.css';

export const CssModuleButton = () => (
  // styles.button === "_2wpxM3yizfwbWee6k0UlD4" at runtime
  <button className={styles.button}>Click me!</button>
);

// index.css (CSS module)
// .button {
//   background-color: #ff0000;
//   width: 320px;
//   padding: 20px;
//   border-radius: 5px;
//   border: none;
//   outline: none;
// }
// .button:hover { color: #fff; }
// .button:active { position: relative; top: 2px; }
// @media (max-width: 480px) { .button { width: 160px; } }

// Dev-time readable class names via css-loader localIdentName
// test: /\.css$/, loader: 'css-loader', options: { modules: { localIdentName: '[local]--[hash:base64:5]' } }

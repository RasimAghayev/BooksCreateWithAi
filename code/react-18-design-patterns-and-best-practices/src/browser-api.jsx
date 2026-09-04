// Event handling: synthetic events, single generic handler (switch), refs, forwardRef, transitions, SVG
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 5 (pp. 111-119)

import { useRef } from 'react';

// Synthetic event + nativeEvent check
export const Button = () => {
  const handleClick = (syntheticEvent) => {
    console.log(syntheticEvent instanceof MouseEvent); // false
    console.log(syntheticEvent.nativeEvent instanceof MouseEvent); // true
  };

  return <button onClick={handleClick}>Click me!</button>;
};

// Single event handler per component (Michael Chan event-switch pattern)
export const DualButton = () => {
  const handleEvent = (event) => {
    switch (event.type) {
      case 'click':
        console.log('clicked');
        break;
      case 'dblclick':
        console.log('double clicked');
        break;
      default:
        console.log('unhandled', event.type);
    }
  };

  return (
    <button onClick={handleEvent} onDoubleClick={handleEvent}>
      Click me!
    </button>
  );
};

// useRef: imperative access to a DOM node (focus the input)
export const FocusInput = () => {
  const inputRef = useRef(null);

  const handleClick = () => inputRef.current.focus();

  return (
    <>
      <input type="text" ref={inputRef} />
      <button onClick={handleClick}>Set Focus</button>
    </>
  );
};

// forwardRef: expose a child DOM node to a parent ref
export const TextInputWithRef = forwardRef((props, ref) => (
  <input ref={ref} type="text" {...props} />
));
TextInputWithRef.displayName = 'TextInputWithRef';

// App consuming TextInputWithRef
export const ForwardRefApp = () => {
  const inputRef = useRef(null);
  const handleClick = () => inputRef.current.focus();

  return (
    <div>
      <TextInputWithRef ref={inputRef} />
      <button onClick={handleClick}>Focus on input</button>
    </div>
  );
};

// CSS-in-JS transition via react-transition-group
export const FadeTransition = () => (
  <TransitionGroup
    transitionName="fade"
    transitionAppear
    transitionAppearTimeout={500}
  >
    <h1>Hello React</h1>
  </TransitionGroup>
);

// SVG as a configurable React component
export const Circle = ({ x, y, radius, fill = 'red' }) => (
  <svg>
    <circle cx={x} cy={y} r={radius} fill={fill} />
  </svg>
);

// Fixed variant of the base circle
export const RedCircle = ({ x, y, radius }) => (
  <Circle x={x} y={y} radius={radius} fill="red" />
);

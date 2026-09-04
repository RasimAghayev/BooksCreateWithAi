// React 18 new APIs & hooks: createRoot, hydrateRoot, renderToPipeableStream,
// useId, useTransition, useDeferredValue, useInsertionEffect, StrictMode
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 10 (pp. 196-211)

import React from 'react';
import { createRoot, hydrateRoot } from 'react-dom/client';

// createRoot — replaces ReactDOM.render; supports concurrent rendering
const App = () => <div>Hello, world!</div>;
const root = createRoot(document.getElementById('root'));
root.render(<App />);

// hydrateRoot — hydrate server-rendered HTML into an interactive tree
// (server & client must render identical markup; use Suspense to align structure)
export const hydrateApp = () => {
  if (root.isMounted()) {
    hydrateRoot(document.getElementById('root'), <App />);
  } else {
    root.render(<App />);
  }
};

// renderToPipeableStream — stream SSR to a Node.js http response
import { renderToPipeableStream } from 'react-dom/server';
import { createServer } from 'http';
export const startStreamServer = () => {
  const server = createServer((req, res) => {
    const stream = renderToPipeableStream(<App />);
    stream.pipe(res);
  });
  server.listen(3000);
};

// useId — SSR-safe unique id (do NOT use for list keys; needs matching SSR/CSR trees)
export const IdExample = () => {
  const id = React.useId();
  return <div id={id}>Hello, world!</div>;
};

// useTransition — mark low-priority updates; isPending shows a loading state
export const TransitionExample = () => {
  const [data, setData] = useState(null);
  const [startTransition, isPending] = React.useTransition();

  const handleClick = () => {
    startTransition(() => {
      const newData = fetchData();
      setData(newData);
    });
  };

  return (
    <div>
      {isPending && <LoadingSpinner />}
      <button onClick={handleClick}>Fetch Data</button>
      {data && <DataDisplay data={data} />}
    </div>
  );
};

// useDeferredValue — defer a value update to the next frame (works WITH useTransition timing)
export const DeferredMove = () => {
  const [x, setX] = useState(0);
  const deferredX = React.useDeferredValue(x);

  function handleClick() {
    setX((v) => v + 100);
  }

  return (
    <div
      style={{ transform: `translateX(${deferredX}px)` }}
      onClick={handleClick}
    >
      Click me!
    </div>
  );
};

// useInsertionEffect — insert DOM nodes before CSS; cleanup removes them
export const InsertionExample = () => {
  React.useInsertionEffect(() => {
    const canvas = document.createElement('canvas');
    canvas.width = 300;
    canvas.height = 200;
    canvas.style.backgroundColor = 'red';
    document.body.appendChild(canvas);
    return () => {
      document.body.removeChild(canvas);
    };
  }, []);

  return (
    <div>
      <h1>Hello, world!</h1>
    </div>
  );
};

// StrictMode — extra dev-time checks (double-invokes components, detects unsafe lifecycles)
export const StrictRoot = () => (
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

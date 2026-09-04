// Next.js universal app: getInitialProps + pages convention + SSR data flow
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 12 (pp. 281-285)

import App from '../pages/api';
import { Provider } from 'react-redux';
import store from '../store';
import AppRouter from '../store';
import App from '../components/App';
import { createRoot } from 'react-dom/client';

// src/pages/index.js — Next.js convention: files in /pages become routes
import fetch from 'isomorphic-fetch';

const App = (props) => {
  return (
    <ul>
      {props.gists.map(gist => (
        <li key={gist.id}>{gist.description}</li>
      ))}
    </ul>
  );
};
export default App;

// Data fetching at the module level (runs on BOTH server & client).
// Returned object becomes `props` on the component.
App.getInitialProps = async () => {
  const url = 'https://api.github.com/users/gaearon/gists';
  const response = await fetch(url);
  const gists = await response.json();
  return { gists };
};

/*
Key points:
- Next.js uses /pages conventions (index.js -> "/"); zero-config SSR.
- getInitialProps runs on server (SSR) and client (client-side nav) → data is universal, no manual dehydration.
- Hot module replacement + fast refresh baked in; `npm run dev` -> http://localhost:3000.
- tsconfig.json (Next uses swc/transpiler; tsconfig still applies to IDE & type checks) — strict, jsx react-jsx, esnext.
*/

// SSR hydration & dehydration pattern (data fetched on server, injected via window, re-rendered on client)
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 12 (pp. 277-282)

import express from 'express';
import fetch from 'isomorphic-fetch';
import { renderToString } from 'react-dom/server';

// Server must accept injected data as props so SSR != empty client render
type Gist = { id: string; description: string };

export const AppSsr = ({ gists }) => (
  <ul>
    {gists.map(gist => (
      <li key={gist.id}>{gist.description}</li>
    ))}
  </ul>
);

// Express route: fetch on server, render string, inject data into HTML
export const makeSsrHandler = () => {
  const app = express();

  app.use(express.static('dist/public'));

  app.get('/', (req, res) => {
    fetch('https://api.github.com/users/gaearon/gists')
      .then(response => response.json())
      .then(gists => {
        const body = renderToString(<AppSsr gists={gists} />);
        const html = `
          <!DOCTYPE html>
          <html>
            <head><meta charset="UTF-8" /></head>
            <body>
              <div id="root">${body}</div>
              <script>window.gists = ${JSON.stringify(gists)}</script>
              <script src="/bundle.js"></script>
            </body>
          </html>`;
        res.send(html);
      });
  });

  return app;
};

/*
Key points:
- renderToString (old API) adds data-reactroot/checksum; client must rehydrate with the SAME props.
- Problem: client render has no `gists` → "Cannot read property 'map' of undefined".
- Fix (dehydration/hydration): inject JSON into a <script> (window.gists), then on client read it back:
    ReactDOM.hydrate(<App gists={window.gists} />, root)
- React 17: hydrateRoot (replaces ReactDOM.hydrate) + createRoot; needs matching trees (suppressHydrationWarning to silence mismatches).
*/

# React 18 Design Patterns and Best Practices — Method

**Book**: React 18 Design Patterns and Best Practices, Fourth Edition
**Author**: Carlos Santana Roldán (Packt Publishing, 2023)
**ISBN**: 978-1-80323-310-9
**Pages**: 525

## Overview

A practical guide to building production-ready React 18 applications using modern design patterns, TypeScript, testing, performance optimization, and deployment best practices.

---

## Chapter 1: Taking Your First Steps with React

**Pages**: 30–41

- **Technical requirements**: Vite as the modern build tool, replacing Create React App
- **Declarative vs Imperative**: React uses declarative programming — describe *what* the UI should look like, not *how* to change it
- **React elements**: `React.createElement()` returns lightweight JavaScript objects that React reconciles against the DOM
- **Virtual DOM**: Efficient diffing algorithm minimizes actual DOM mutations
- **JavaScript fatigue**: Vite solves startup/build time issues; use native ES modules in development
- **Key commands**:
  - `npm create vite@latest my-app` — scaffold with Vite
  - `npm run dev` / `npm run build` — dev server / production build

---

## Chapter 2: Introducing TypeScript

**Pages**: 44–57

- **TypeScript's features**: Static typing, interfaces, enums, namespaces, template literals
- **Types**: `string`, `number`, `boolean`, `array`, `tuple`, `enum`, `any`, `unknown`, `never`
- **Interfaces**: Define object shapes; `interface` vs `type`; extending interfaces
- **Enums**: `enum Direction { Up, Down }`
- **Namespaces**: Group related code (prefer modules over namespaces in modern TS)
- **Template literals**: `Hello ${name}` with type safety
- **TS config**: `tsconfig.json` with `strict`, `noImplicitAny`, `strictNullChecks`

---

## Chapter 3: Cleaning Up Your Code

**Pages**: 58–85

- **JSX**: Syntactic sugar for `React.createElement()`
- **Props and children**: Component composition via `props.children`
- **Babel**: Compiles JSX/TypeScript to browser-compatible JavaScript
- **Functional programming**:
  - First-class functions (functions as values)
  - Pure functions (no side effects)
  - Immutability (`Object.freeze`, spread operator)
  - Currying: `fn(a)(b)` instead of `fn(a, b)`
  - Composition: passing components as children or props
- **Styling tools**: EditorConfig, Prettier, ESLint
- **Git Hooks**: Husky + lint-staged for pre-commit checks

---

## Chapter 4: Exploring Popular Composition Patterns

**Pages**: 88–101

- **children prop pattern**: Pass JSX as children for flexible composition
- **Container and presentational patterns**: Separate data-fetching (container) from UI rendering (presentational)
- **Higher-Order Components (HOCs)**: Functions that take a component and return a new component with enhanced behavior
- **FunctionAsChild (render props)**: Pass a function as a child for dynamic rendering

---

## Chapter 5: Writing Code for the Browser

**Pages**: 102–120

- **Forms**: Uncontrolled components (ref-based) vs controlled components (state-based)
- **Events**: Synthetic events, event handlers (`onClick`, `onChange`, `onSubmit`)
- **Refs**: `useRef` for accessing DOM nodes; `forwardRef` for forwarding refs to child components
- **Animations**: CSS transitions/animations via React state
- **SVG**: Inline SVG rendering in React

---

## Chapter 6: Making Your Components Look Beautiful

**Pages**: 122–147

- **CSS in JavaScript**: Dynamic styling via libraries like styled-components
- **Inline styles**: Direct style objects (`style={{ color: 'red' }}`)
- **CSS Modules**: Locally scoped CSS via Webpack configuration
- **styled-components**: Tagged template literals for component-level styling
- **Webpack 5**: CSS loader, PostCSS, and module configuration for scoped styles

---

## Chapter 7: Anti-Patterns to Be Avoided

**Pages**: 150–157

- **Initializing state from props**: Use `null`/default values instead of copying props into state
- **Using array indexes as keys**: Use stable, unique IDs for list items
- **Spreading properties on DOM elements**: Can leak unknown attributes to DOM nodes (e.g., `{...props}`)

---

## Chapter 8: React Hooks

**Pages**: 160–198

- **useState**: State management in functional components
- **Rules of Hooks**: Only call hooks at top level; only from React functions
- **Migrating class components**: Replace `this.state`/`this.setState` with `useState`, `componentDidMount` with `useEffect`
- **useEffect**: Side effects with cleanup; dependency array optimization
- **useCallback**: Memoize function references
- **useMemo**: Memoize computed values
- **memo**: Prevent unnecessary re-renders of child components
- **useReducer**: Alternative to useState for complex state logic; Redux-like reducer pattern

---

## Chapter 9: React Router

**Pages**: 200–222

- **v6.4+**: Data loading with `loader` functions, `Route` as JSX
- **Routes and Route**: Nested routing structure
- **URL parameters**: Dynamic route params via `useParams`
- **Programmatic navigation**: `useNavigate`
- **Link vs NavLink**: Active link styling
- **Router loaders**: Fetch data before rendering a route

---

## Chapter 10: React 18 New Features

**Pages**: 224–241

- **Concurrent mode**: Interruptible rendering for better responsiveness
- **Automatic batching**: Multiple state updates batch automatically
- **Transitions**: `useTransition()` for non-urgent updates with `isPending`
- **Suspense on the server**: Stream server-rendered content
- **New APIs**:
  - `createRoot()` — replaces `ReactDOM.render()`
  - `hydrateRoot()` — replaces `ReactDOM.hydrate()`
  - `renderToPipeableStream()` / `renderToReadableStream()` — streaming SSR
- **New Hooks**:
  - `useId` — SSR-safe unique ID generation
  - `useTransition` — `startTransition`, `isPending`
  - `useDeferredValue` — defer value updates
  - `useInsertionEffect` — inject CSS before DOM mutations
- **Strict mode**: Double-invokes render and effects in dev for catching bugs
- **Node 18+**: Required for latest React 18 SSR features

---

## Chapter 11: Managing Data

**Pages**: 242–262

- **React Context API**: `createContext`, `useContext`, `Context.Provider` for global state
- **React Suspense with SWR**: Data fetching with `<Suspense>` fallback and `useSWR` hook
- **Pokedex demo**: Real project using Pokemon API, Suspense for loading states
- **Redux Toolkit**:
  - `configureStore` — opinionated store setup
  - `createSlice` — slices with reducers and actions
  - `combineReducers` / `combineSlices` — combine multiple reducers
  - `useSelector` / `useDispatch` — connect React components to Redux store

---

## Chapter 12: Server-Side Rendering

**Pages**: 264–283

- **Universal applications**: Same code runs on server and client
- **SSR benefits**: SEO, better initial load performance, social sharing previews
- **Webpack configuration**: Multi-config setup (server + client bundles)
- **Manual SSR**: `renderToString` from `react-dom/server`
- **Next.js**: File-based routing, automatic SSR, static generation (`getStaticProps`, `getServerSideProps`)

---

## Chapter 13: Understanding GraphQL with a Real Project

**Pages**: 286–358

- **Backend stack**: Apollo Server, GraphQL, Sequelize ORM, PostgreSQL, JWT authentication
- **GraphQL schema**: Scalar types, queries, mutations, type merging
- **Resolvers**: Query and mutation resolvers with Sequelize models
- **Authentication**: JWT token generation (`jwt.sign`), verification (`jwt.verify`), password hashing
- **Frontend**: Apollo Client, GraphQL queries/mutations with `useQuery`/`useMutation`
- **Real project**: Full login system with user registration, authentication, and dashboard

---

## Chapter 14: MonoRepo Architecture

**Pages**: 360–439

- **NPM Workspaces**: Shared dependencies across packages
- **Packages**: utils, API, frontend, devtools
- **devtools package**: Webpack configs (common, dev, prod), shared build tooling
- **utils package**: Shared utility functions
- **API package**: Shared GraphQL types, resolvers, Sequelize models
- **Frontend package**: React app with routing, login system, site configuration
- **Webpack configs**: Common (Babel, CSS), dev (HMR, dev server), prod (minification, tree-shaking)

---

## Chapter 15: Improving the Performance of Your Applications

**Pages**: 440–447

- **Reconciliation**: React's diffing algorithm uses `key` prop to identify elements
- **Keys**: Use stable, unique string IDs; avoid array indexes
- **Optimization techniques**:
  - Code splitting: dynamic `import()`
  - Lazy loading: `React.lazy` + `Suspense`
  - Memoization: `React.memo`, `useMemo`, `useCallback`
- **Tools and libraries**:
  - Bundle analyzers (`webpack-bundle-analyzer`)
  - Image optimization (lazy loading, WebP)
  - Babel plugins for optimization (React transform, dead code elimination)

---

## Chapter 16: Testing and Debugging

**Pages**: 450–472

- **Jest**: Test runner with snapshot testing, mocking
- **React Testing Library**: Render components, find elements, simulate events
- **Vitest**: Vite-native test runner; drop-in Jest replacement
- **Vitest configuration**: `vitest.config.ts`, globals, environment
- **In-source testing**: Co-locate tests with components
- **React DevTools**: Component tree inspection, props/state inspection
- **Redux DevTools**: Time-travel debugging, action/state inspection

---

## Chapter 17: Deploying to Production

**Pages**: 474–505

- **DigitalOcean Droplet**: VPS setup with Node.js
- **NGINX**: Reverse proxy, static file serving, SSL termination
- **PM2**: Process manager for keeping Node.js apps running
- **Domain configuration**: DNS, SSL certificates
- **CI/CD**: CircleCI pipeline with SSH key authentication, environment variables, automated builds/deploy

---

## Key Design Patterns

| Pattern | Description | Chapter |
|---------|-------------|---------|
| Container/Presentational | Separate data logic from UI rendering | Ch 4 |
| Higher-Order Component (HOC) | Wrap components with reusable logic | Ch 4 |
| Function-as-Child | Render props via functions | Ch 4 |
| children prop | Pass JSX as component children | Ch 4 |
| Compound | Related components share implicit state | Ch 8 |
| Custom Hooks | Reusable stateful logic (e.g., `useFetch`) | Ch 8 |

## Key React 18 APIs

| API | Purpose |
|-----|---------|
| `createRoot()` | Modern root rendering |
| `hydrateRoot()` | SSR hydration |
| `useId()` | SSR-safe unique IDs |
| `useTransition()` | Non-blocking state updates |
| `useDeferredValue()` | Defer value propagation |
| `useInsertionEffect()` | CSS injection before DOM mutations |
// Enables `import styles from './index.css'` with TypeScript.
// styles.button -> scoped (hashed) class name at runtime.
declare module '*.css' {
  const content: Record<string, string>;
  export default content;
}

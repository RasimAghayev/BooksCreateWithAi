// First-class functions / Higher-Order Functions (HoFs) in React/FP
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 3 (pp. 82-87)

// HoF: takes a function, returns an enhanced function
const log = fn => (...args) => {
  console.log('args:', args);
  return fn(...args);
};

const add = (x, y) => x + y;
const logAdd = log(add);
// logAdd(2, 3) -> logs [2,3] then returns 5

export { add, log, logAdd };

// Purity: same input -> same output, no side effects
const pureAdd = (x, y) => x + y;

// Impure: mutates outer scope -> different result on repeat calls
let x = 0;
const impureAdd = y => (x = x + y);

export { pureAdd, impureAdd };

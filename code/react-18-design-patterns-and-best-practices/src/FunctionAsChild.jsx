// FunctionAsChild pattern: children is a function that receives params from the parent
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 4 (pp. 100-105)

import React from 'react';

// Parent calls the children function inside its render
const Name = ({ children }) => children('World');
export default Name;

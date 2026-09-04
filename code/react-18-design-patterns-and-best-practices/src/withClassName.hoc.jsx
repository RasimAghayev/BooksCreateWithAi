// withClassName HOC: attaches a className to every wrapped component
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 4 (p. 97)

import React from 'react';

const withClassName = Component => props => (
  <Component {...props} className="my-class" />
);

export default withClassName;

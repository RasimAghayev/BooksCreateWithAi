import React from 'react';

const Hello = ({ name }) => (
  <h1 className="Hello">Hello {name || 'World'}</h1>
);

Hello.defaultProps = { name: '' };

export default Hello;

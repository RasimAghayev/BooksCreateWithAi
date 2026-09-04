import React from 'react';
import { render, cleanup } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import Hello from './index';

describe('Hello Component', () => {
  it('should render Hello World', () => {
    const { getByText } = render(<Hello />);
    expect(getByText('Hello World')).toBeInTheDocument();
  });

  it('should render the name prop', () => {
    const { getByText } = render(<Hello name="Carlos" />);
    expect(getByText('Hello Carlos')).toBeInTheDocument();
  });

  it('should has .Hello classname', () => {
    const { container } = render(<Hello />);
    expect(container.firstChild).toHaveClass('Hello');
  });

  afterAll(cleanup);
});

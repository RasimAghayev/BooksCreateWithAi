import { cleanup, render } from '@testing-library/react';
import { afterAll, describe, expect, it } from 'vitest';
import Hello from './index';

describe('Hello Component', () => {
  it('should render Hello World', () => {
    const wrapper = render(<Hello />);
    expect(wrapper.getByText('Hello World')).toBeDefined();
  });

  it('should render the name prop', () => {
    const wrapper = render(<Hello name="Carlos" />);
    expect(wrapper.getByText('Hello Carlos')).toBeDefined();
  });

  it('should has .Hello classname', () => {
    const wrapper = render(<Hello />);
    const firstChild = wrapper.container.firstChild as HTMLElement;
    expect(firstChild?.classList.contains('Hello')).toBe(true);
  });

  afterAll(cleanup);
});

import React from 'react';
import { render, cleanup, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import ShowInformation from './index';

describe('Show Information Component', () => {
  let wrapper;

  beforeEach(() => {
    wrapper = render(<ShowInformation />);
  });

  it('should modify the name', () => {
    const nameInput = wrapper.container.querySelector('input[name="name"]');
    const ageInput = wrapper.container.querySelector('input[name="age"]');
    fireEvent.change(nameInput, { target: { value: 'Carlos' } });
    fireEvent.change(ageInput, { target: { value: 34 } });
    expect(nameInput.value).toBe('Carlos');
    expect(ageInput.value).toBe('34');
  });

  it('should show the personal information when user clicks on the button', () => {
    const button = wrapper.container.querySelector('button');
    fireEvent.click(button);
    const showInformation = wrapper.container.querySelector('.personalInformation');
    expect(showInformation).toBeInTheDocument();
  });

  afterAll(cleanup);
});

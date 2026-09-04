// Forms: uncontrolled → single-handler → controlled pattern
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 5 (pp. 103-108)

import { useState } from 'react';

// Uncontrolled: DOM holds the state; read via ref on submit
export const UncontrolledForm = () => {
  const [value, setValue] = useState('');

  const handleChange = (e) => setValue(e.target.value);
  const handleSubmit = (e) => {
    e.preventDefault();
    console.log(value);
  };

  return (
    <form onSubmit={handleSubmit}>
      <input type="text" onChange={handleChange} />
      <button>Submit</button>
    </form>
  );
};

// Multi-field with ONE handler keyed by input name
export const UncontrolledMultiField = () => {
  const [values, setValues] = useState({ firstName: '', lastName: '' });

  const handleChange = ({ target: { name, value } }) =>
    setValues({ ...values, [name]: value });

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log(`${values.firstName} ${values.lastName}`);
  };

  return (
    <form onSubmit={handleSubmit}>
      <input name="firstName" type="text" onChange={handleChange} />
      <input name="lastName" type="text" onChange={handleChange} />
      <button>Submit</button>
    </form>
  );
};

// Controlled: React state drives the field value
export const ControlledForm = () => {
  const [values, setValues] = useState({ firstName: '', lastName: '' });

  const handleChange = ({ target: { name, value } }) =>
    setValues({ ...values, [name]: value });

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log(`${values.firstName} ${values.lastName}`);
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        name="firstName"
        type="text"
        value={values.firstName}
        onChange={handleChange}
      />
      <input
        name="lastName"
        type="text"
        value={values.lastName}
        onChange={handleChange}
      />
      <button>Submit</button>
    </form>
  );
};

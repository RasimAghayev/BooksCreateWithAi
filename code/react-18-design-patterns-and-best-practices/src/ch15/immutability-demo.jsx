import React, { useState } from 'react';
import { render } from 'react-dom';

const ImmutabilityDemo = () => {
  const [name, setName] = useState('');
  const [age, setAge] = useState(0);
  const [show, setShow] = useState(false);

  const handleOnChange = (e) => {
    const { name, value } = e.target;
    if (name === 'name') setName(value);
    if (name === 'age') setAge(Number(value));
  };

  const handleShowInformation = () => setShow(true);

  if (show) {
    return (
      <div className="ShowInformation">
        <h1>Personal Information</h1>
        <div className="personalInformation">
          <p><strong>Name:</strong> {name}</p>
          <p><strong>Age:</strong> {age}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="ShowInformation">
      <h1>Personal Information</h1>
      <p><strong>Name:</strong></p>
      <p>
        <input name="name" type="text" value={name} onChange={handleOnChange} />
      </p>
      <p>
        <input name="age" type="number" value={age} onChange={handleOnChange} />
      </p>
      <p><button onClick={handleShowInformation}>Show Information</button></p>
    </div>
  );
};

export default ImmutabilityDemo;

import React from 'react';
import PropTypes from 'prop-types';

const ReconciliationDemo = ({ addAtStart }) => {
  const [items, setItems] = React.useState([
    { id: 1, name: 'Alice' },
    { id: 2, name: 'Bob' },
    { id: 3, name: 'Charlie' },
  ]);

  const addItem = () => {
    const newId = Math.max(...items.map((i) => i.id)) + 1;
    setItems(addAtStart
      ? [{ id: newId, name: `User ${newId}` }, ...items]
      : [...items, { id: newId, name: `User ${newId}` }]
    );
  };

  return (
    <div className="reconciliation-demo">
      <button onClick={addItem}>Add Item</button>
      <ul>
        {items.map((item) => (
          <li key={item.id}>{item.name}</li>
        ))}
      </ul>
    </div>
  );
};

ReconciliationDemo.propTypes = {
  addAtStart: PropTypes.bool,
};

ReconciliationDemo.defaultProps = {
  addAtStart: false,
};

export default ReconciliationDemo;

// React memoization: memo (component), useMemo (value), useCallback (function), useReducer
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 8 (pp. 175-190)

import {
  useState,
  useEffect,
  memo,
  useMemo,
  useCallback,
  useReducer,
} from 'react';

// --- memo: shallow-compare props, skip re-render ---
const TaskBase = ({ id, task, handleDelete }) => {
  useEffect(() => console.log('Rendering <Task />', task));
  return (
    <li>
      {task} <button onClick={() => handleDelete(id)}>X</button>
    </li>
  );
};
export const Task = memo(TaskBase);

// --- useMemo: memoize an expensive computed value w/ change-tracking deps ---
export const TodoList = ({ todoList, term }) => {
  const filteredTodoList = useMemo(
    () => todoList.filter(t => t.task.toLowerCase().includes(term.toLowerCase())),
    [term, todoList] // <-- must list all referenced values or the hook is stale
  );

  return (
    <ul>
      {filteredTodoList.map(t => (
        <li key={t.id}>{t.task}</li>
      ))}
    </ul>
  );
};

// --- useCallback: memoize a function identity (for props/effects deps) ---
export const DeleteApp = ({ todoList, setTodoList }) => {
  // Stable reference so memoized children/effects don't re-render on every parent render
  const handleDelete = useCallback(
    (taskId) => {
      setTodoList(todoList.filter(t => t.id !== taskId));
    },
    [todoList] // re-created only when the dependency changes
  );

  // When a memoized fn is passed into useEffect, it must be a dep of the effect
  const printTodoList = useCallback(() => {
    console.log('Changing todoList', todoList);
  }, [todoList]);

  useEffect(() => {
    printTodoList();
  }, [printTodoList]);

  return (
    <ul>
      {todoList.map(t => (
        <li key={t.id}>{t.task}</li>
      ))}
    </ul>
  );
};

// --- useReducer: Redux-like dispatch/reducer for complex state logic ---
// reducer(state, action) -> newState
const todosReducer = (state, action) => {
  switch (action.type) {
    case 'add':
      return [...state, action.payload];
    case 'delete':
      return state.filter(t => t.id !== action.payload.id);
    case 'update':
      return state.map(t =>
        t.id === action.payload.id ? action.payload : t
      );
    default:
      return state;
  }
};

export const Notes = () => {
  const [notes, dispatch] = useReducer(todosReducer, []);

  const addNote = (note) => dispatch({ type: 'add', payload: note });
  const deleteNote = (id) => dispatch({ type: 'delete', payload: { id } });
  const updateNote = (note) => dispatch({ type: 'update', payload: note });

  // dispatch only takes plain-object actions (no middleware/thunk/saga built-in)
  useEffect(() => {
    console.log('Notes:', notes);
  }, [notes]);

  return (
    <button onClick={() => addNote({ id: Date.now(), task: 'New' })}>+</button>
  );
};

export default Notes;

/*
Recap cheat sheet:
- memo(Component)            -> shallow-compares props; skip re-render if unchanged.
- useMemo(fn, deps)          -> memoizes a VALUE (computed property / heavy calc).
- useCallback(fn, deps)      -> memoizes a FUNCTION IDENTITY (same as useMemo(() => fn, deps)).
- useReducer(reducer, init)  -> dispatch(action) -> reducer(state, action) -> state. Redux-like, no store/middleware.

Golden rule: don't use memo/useMemo/useCallback until you actually see a perf problem.
*/

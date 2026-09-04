const reduxDevtools = require('redux-devtools-extension');
const { createStore, applyMiddleware } = require('redux');
const thunk = require('redux-thunk').default;
const { composeWithDevTools } = reduxDevtools;
const rootReducer = require('@reducers');

function configureStore({ initialState, reducer }) {
  const middleware = [thunk];
  return createStore(
    rootReducer,
    initialState,
    composeWithDevTools(applyMiddleware(...middleware))
  );
}

module.exports = configureStore;

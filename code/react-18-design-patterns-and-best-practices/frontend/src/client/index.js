const { ApolloClient, ApolloProvider, InMemoryCache } = require('@apollo/client');
const { render } = require('react-dom');
const config = require('../config');
const AppRoutes = require('./AppRoutes');

const client = new ApolloClient({
  uri: config.api.uri,
  cache: new InMemoryCache(),
});

render(
  <ApolloProvider client={client}>
    <AppRoutes />
  </ApolloProvider>,
  document.querySelector('#root')
);

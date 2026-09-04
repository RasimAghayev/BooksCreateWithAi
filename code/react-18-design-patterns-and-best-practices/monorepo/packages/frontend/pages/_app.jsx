const { ApolloClient, InMemoryCache } = require('@apollo/client');
const React = require('react');
const Config = require('@/config');
const GlobalStyle = require('@/components/GlobalStyles/GlobalStyles');
const { useApollo } = require('@/contexts/apolloClient');
const FormProvider = require('@/contexts/form');
const UserProvider = require('@/contexts/user');

const App = ({ Component, pageProps }) => {
  const apolloClient = useApollo((pageProps && pageProps.initialApolloState) || {});
  return React.createElement(React.Fragment, null,
    React.createElement(GlobalStyle, null),
    React.createElement(require('@apollo/client').ApolloProvider, { client: apolloClient },
      React.createElement(UserProvider, null,
        React.createElement(FormProvider, null,
          React.createElement(Component, { ...pageProps })
        )
      )
    )
  );
};

App.getInitialProps = async () => ({
  ...Config,
});

module.exports = App;

const React = require('react');
const { createContext, useContext, useState, useEffect, useMemo } = React;
const { useMutation, useQuery } = require('@apollo/client');
const { useCookies } = require('react-cookie');
const { getGraphQlError, redirectTo } = require('@web-creator/utils');
const Config = require('@/config');
const GET_USER_QUERY = require('@/graphql/user/getUser.query');
const LOGIN_MUTATION = require('@/graphql/user/login.mutation');

const UserContext = createContext({
  login: () => null,
  user: null,
});

const UserProvider = ({ children }) => {
  const [cookies, setCookie] = useState({});
  const [user, setUser] = useState(null);
  const [loginMutation] = useMutation(LOGIN_MUTATION);
  const { data: dataUser } = useQuery(GET_USER_QUERY, {
    variables: { at: cookies[`at-${Config.site}`] || '' },
  });

  useEffect(() => {
    if (dataUser) {
      setUser(dataUser.getUser);
    }
  }, [dataUser]);

  async function login(input) {
    try {
      const { data: dataLogin } = await loginMutation({
        variables: {
          emailOrUsername: input.emailOrUsername,
          password: input.password,
        },
      });
      if (dataLogin) {
        setCookie(`at-${Config.site}`, dataLogin.login.token, {
          path: '/',
          maxAge: 45 * 60 * 1000,
        });
        return dataLogin.login.token;
      }
    } catch (err) {
      return getGraphQlError(err);
    }
    return null;
  }

  const context = useMemo(() => ({ login, user }), [user]);
  return React.createElement(UserContext.Provider, { value: context }, children);
};

module.exports = { UserContext, UserProvider };

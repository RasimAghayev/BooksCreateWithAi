const { createContext, useContext, useState, useEffect } = require('react');
const { useMutation, useQuery } = require('@apollo/client');
const { useCookies } = require('react-cookie');
const GET_USER_QUERY = require('../graphql/user/getUser.query');
const LOGIN_MUTATION = require('../graphql/user/login.mutation');

const UserContext = createContext({
  login: () => null,
  connectedUser: null,
});

const UserProvider = ({ page = '', children }) => {
  const [cookies, setCookie] = useCookies();
  const [connectedUser, setConnectedUser] = useState(null);
  const [loginMutation] = useMutation(LOGIN_MUTATION);
  const { data: dataUser } = useQuery(GET_USER_QUERY, {
    variables: { at: cookies.at || '' },
  });

  useEffect(() => {
    if (dataUser) {
      if (!dataUser.getUser.id && page !== 'login') {
        window.location.href = '/login?redirectTo=/dashboard';
      } else {
        setConnectedUser(dataUser.getUser);
      }
    }
  }, [dataUser, page]);

  async function login(input) {
    try {
      const { data: dataLogin } = await loginMutation({
        variables: { email: input.email, password: input.password },
      });
      if (dataLogin) {
        setCookie('at', dataLogin.login.token, { path: '/' });
        return dataLogin.login.token;
      }
    } catch (err) {
      console.error(err);
      return null;
    }
  }

  const context = { login, connectedUser };
  return <UserContext.Provider value={context}>{children}</UserContext.Provider>;
};

module.exports = { UserContext, UserProvider };

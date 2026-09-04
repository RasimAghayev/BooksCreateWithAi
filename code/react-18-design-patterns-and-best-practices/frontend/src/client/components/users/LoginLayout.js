const { useContext } = require('react');
const { UserContext } = require('../../contexts/user');
const Login = require('./Login');

const LoginLayout = ({ currentUrl }) => {
  const { login } = useContext(UserContext);
  return <Login login={login} currentUrl={currentUrl} />;
};

module.exports = LoginLayout;

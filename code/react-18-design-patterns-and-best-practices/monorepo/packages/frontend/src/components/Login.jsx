const React = require('react');
const { useContext, useState } = React;
const { UserContext } = require('@/contexts/user');
const { FormContext } = require('@/contexts/form');
const { redirectTo, getRedirectToUrl } = require('@web-creator/utils');
const { Button, Input, RenderIf } = require('@web-creator/design-system');
const { CSS } = require('./Login.styled');

const Login = () => {
  const redirectToUrl = getRedirectToUrl();
  const [values, setValues] = useState({ emailOrUsername: '', password: '' });
  const [notification, setNotification] = useState({ id: Math.random(), message: '' });
  const [invalidLogin, setInvalidLogin] = useState(false);

  const { change } = useContext(FormContext);
  const { login } = useContext(UserContext);

  const onChange = (e) => change(e, setValues);

  const handleSubmit = async (user) => {
    const response = await login(user);
    if (response && response.error) {
      setInvalidLogin(true);
      setNotification({ id: Math.random(), message: response.message });
    } else {
      redirectTo(redirectToUrl || '/', true);
    }
  };

  return React.createElement(React.Fragment, null,
    React.createElement(RenderIf, { isTrue: invalidLogin && notification.message !== '' },
      notification.message
    ),
    React.createElement(CSS.Login, null,
      React.createElement('header', null,
        React.createElement('img', { className: 'logo', src: '/images/isotype.png', alt: 'Logo' }),
        React.createElement('br'),
        React.createElement('h2', null, 'Sign In')
      ),
      React.createElement('section', null,
        React.createElement(Input, {
          autoComplete: 'off',
          name: 'emailOrUsername',
          placeholder: 'Email Or Username',
          onChange: onChange,
          value: values.emailOrUsername,
        }),
        React.createElement(Input, {
          name: 'password',
          type: 'password',
          placeholder: 'Password',
          onChange: onChange,
          value: values.password,
        }),
        React.createElement('div', { className: 'actions' },
          React.createElement(Button, { onClick: () => handleSubmit(values) }, 'Login'),
          React.createElement(Button, { color: 'success' }, 'Register')
        )
      )
    )
  );
};

module.exports = Login;

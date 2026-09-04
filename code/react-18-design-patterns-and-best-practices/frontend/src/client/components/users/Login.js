const { useState } = require('react');
const StyledLogin = require('./Login.styled');

const Login = ({ login, currentUrl }) => {
  const [values, setValues] = useState({
    email: '',
    password: '',
  });
  const [errorMessage, setErrorMessage] = useState('');
  const [invalidLogin, setInvalidLogin] = useState(false);

  const onChange = (e) => {
    const { name, value } = e.target;
    if (name) {
      setValues((prevValues) => ({
        ...prevValues,
        [name]: value,
      }));
    }
  };

  const handleSubmit = async () => {
    const response = await login(values);
    if (response && response.error) {
      setInvalidLogin(true);
      setErrorMessage(response.message);
    } else {
      window.location.href = currentUrl || '/';
    }
  };

  return (
    <StyledLogin>
      <div className="wrapper">
        {invalidLogin && <div className="alert">{errorMessage}</div>}
        <div className="form">
          <p>
            <input
              autoComplete="off"
              type="email"
              className="email"
              name="email"
              placeholder="Email"
              onChange={onChange}
              value={values.email}
            />
          </p>
          <p>
            <input
              autoComplete="off"
              type="password"
              className="password"
              name="password"
              placeholder="Password"
              onChange={onChange}
              value={values.password}
            />
          </p>
          <div className="actions">
            <button name="login" onClick={handleSubmit}>
              Login
            </button>
          </div>
        </div>
      </div>
    </StyledLogin>
  );
};

module.exports = Login;

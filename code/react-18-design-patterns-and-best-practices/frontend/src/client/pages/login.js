const { useState, useContext, useEffect } = require('react');
const { UserContext } = require('../contexts/user');
const LoginLayout = require('../components/users/LoginLayout');

const LoginPage = ({ currentUrl = '' }) => (
  <UserProvider page="login">
    <LoginLayout currentUrl={currentUrl} />
  </UserProvider>
);

module.exports = LoginPage;

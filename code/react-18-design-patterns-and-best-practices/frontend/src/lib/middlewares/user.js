const jwt = require('jsonwebtoken');
const config = require('../../config');

const isAuthenticated = async (token) => {
  if (!token) return false;
  try {
    const decoded = jwt.verify(token, config.security.secretKey);
    return !!decoded;
  } catch (err) {
    return false;
  }
};

const isConnected = (required) => (req, res, next) => {
  const token = req.cookies.at || req.headers.authorization?.split(' ')[1];
  isAuthenticated(token).then((connected) => {
    if (required && !connected) {
      return res.redirect('/login');
    }
    if (!required && connected) {
      return res.redirect('/dashboard');
    }
    req.userConnected = connected;
    next();
  });
};

module.exports = { isConnected, isAuthenticated };

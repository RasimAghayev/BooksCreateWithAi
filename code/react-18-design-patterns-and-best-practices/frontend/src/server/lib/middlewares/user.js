const jwt = require('jsonwebtoken');
const { getUserData } = require('../jwt');

const isConnected = (isLogged = true, roles = ['user'], redirectTo = '/') =>
  async (req, res, next) => {
    const user = await getUserData(req.cookies.at);
    if (!user && !isLogged) {
      return next();
    }
    if (user && isLogged) {
      if (roles.includes('god') && user.role === 'god') {
        return next();
      }
      if (roles.includes('admin') && user.role === 'admin') {
        return next();
      }
      if (roles.includes('user') && user.role === 'user') {
        return next();
      }
      return res.redirect(redirectTo);
    } else {
      return res.redirect(redirectTo);
    }
  };

module.exports = { isConnected };

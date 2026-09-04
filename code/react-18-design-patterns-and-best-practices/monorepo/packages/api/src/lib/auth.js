const jwt = require('jsonwebtoken');
const bcrypt = require('bcryptjs');
const { Buffer } = require('buffer');

const SECRET_KEY = process.env.SECURITY_SECRET_KEY || 'default_secret_key';

const encrypt = (password) => {
  return bcrypt.hashSync(password, 10);
};

const jwtVerify = (accessToken, cb) => {
  jwt.verify(accessToken, SECRET_KEY, (error, accessTokenData = {}) => {
    const { data: user } = accessTokenData;
    if (error || !user) {
      return cb(null);
    }
    const userData = JSON.parse(Buffer.from(user, 'base64').toString('utf-8'));
    return cb(userData);
  });
};

const getUserData = async (accessToken) => {
  const userPromise = new Promise((resolve) =>
    jwtVerify(accessToken, (user) => resolve(user))
  );
  const user = await userPromise;
  return user;
};

const getUserBy = async (criteria, roles, models) => {
  try {
    const user = await models.User.findOne({ where: criteria });
    if (user && roles.includes(user.role)) {
      return user;
    }
    return null;
  } catch (err) {
    return null;
  }
};

const authenticate = async (emailOrUsername, password, models) => {
  try {
    const user = await models.User.findOne({
      where: {
        [models.Sequelize.Op.or]: [
          { email: emailOrUsername },
          { username: emailOrUsername },
        ],
      },
    });
    if (!user) {
      return { error: 'Invalid Login', message: 'Invalid Login' };
    }
    const isPasswordMatch = bcrypt.compareSync(password, user.password);
    if (!isPasswordMatch) {
      return { error: 'Invalid Login', message: 'Invalid Login' };
    }
    if (!user.active) {
      return { error: 'Account inactive', message: 'Your account is not activated yet' };
    }
    const payload = Buffer.from(JSON.stringify(user)).toString('base64');
    const token = jwt.sign({ data: payload }, SECRET_KEY, { expiresIn: '7d' });
    return { token };
  } catch (err) {
    return { error: 'Login failed', message: 'Login failed' };
  }
};

const createToken = (user) => {
  const payload = Buffer.from(JSON.stringify(user)).toString('base64');
  return jwt.sign({ data: payload }, SECRET_KEY, { expiresIn: '7d' });
};

module.exports = {
  encrypt,
  jwtVerify,
  getUserData,
  getUserBy,
  authenticate,
  createToken,
  SECRET_KEY,
};

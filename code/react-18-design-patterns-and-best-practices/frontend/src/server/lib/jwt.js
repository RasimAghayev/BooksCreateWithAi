const jwt = require('jsonwebtoken');
const config = require('../../config');

const { security: { secretKey } } = config;

function jwtVerify(accessToken, cb) {
  jwt.verify(accessToken, secretKey, (error, accessTokenData = {}) => {
    const { data: user } = accessTokenData;
    if (error || !user) {
      return cb(null);
    }
    const userData = Buffer.from(user, 'base64').toString('utf-8');
    return cb(JSON.parse(userData));
  });
}

async function getUserData(accessToken) {
  const userPromise = new Promise((resolve) =>
    jwtVerify(accessToken, (user) => resolve(user))
  );
  const user = await userPromise;
  return user;
}

module.exports = { jwtVerify, getUserData };

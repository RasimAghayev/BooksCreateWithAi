// Sequelize DB connection (PostgreSQL) + JWT auth helpers
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 13 (pp. 305-308)

import { Sequelize } from 'sequelize';
import jwt from 'jsonwebtoken';
import { encrypt, isPasswordMatch, getBase64, setBase64 } from '@contentpi/lib';
import { $db, $security } from '../config';

// --- DB connection: build URI from $db env-backed config ---
const { dialect, port, host, database, username, password } = $db;
const uri = `${dialect}://${username}:${password}@${host}:${port}/${database}`;
const sequelize = new Sequelize(uri);

// models/index.js: attach User to the registry
// const models = { User: require('./User').default(sequelize, Sequelize), sequelize };
// WARNING: { alter: true, force: false } — force drops tables/data!

// --- JWT helpers ---
const { secretKey, expiresIn } = $security;

export const jwtVerify = (accessToken, cb) => {
  jwt.verify(accessToken, secretKey, (error, accessTokenData = {}) => {
    const { data: user } = accessTokenData;
    if (error || !user) return cb(false);
    // Token payload stored as Base64-encoded JSON
    const userData = getBase64(user);
    return cb(userData);
  });
};

export const getUserData = async (accessToken) => {
  const promise = new Promise((resolve) =>
    jwtVerify(accessToken, (user) => resolve(user))
  );
  return promise; // user or false
};

export const createToken = async (user) => {
  const { id, username, password, email, role, active } = user;
  // "token" is an alias for the (already-encryped) password
  const token = setBase64(`${encrypt($security.secretKey)}${password}`);
  const userData = { id, username, email, role, active, token };
  // Sign the JWT: payload = base64(user), signed w/ secret, expiresIn
  const _createToken = jwt.sign(
    { data: setBase64(userData) },
    secretKey,
    { expiresIn }
  );
  return Promise.all([_createToken]);
};

// --- Auth functions (login flow) ---
export const getUserBy = async (where, models) =>
  models.User.findOne({ where, raw: true });

export const doLogin = async (email, password, models) => {
  const user = await getUserBy({ email }, models);
  if (!user) throw new Error('Invalid Login');

  const passwordMatch = isPasswordMatch(encrypt(password), user.password);
  if (!passwordMatch) throw new Error('Invalid Login');

  const isActive = user.active;
  if (!isActive) throw new Error('Your account is not activated yet');

  const [token] = await createToken(user);
  return { token };
};

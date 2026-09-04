module.exports = {
  isProduction: process.env.NODE_ENV === 'production',
  isDevelopment: process.env.NODE_ENV !== 'production',

  PORT: Number(process.env.PORT) || 3000,
  DEV_SERVER_PORT: 3001,
  GRAPHQL_PORT: 4000,
  GRAPHQL_SERVER: process.env.NODE_ENV !== 'production' ? 'localhost' : 'localhost',

  domain: 'localhost',
  baseUrl: process.env.NODE_ENV === 'production'
    ? `https://localhost:3000`
    : `http://localhost:3000`,
  publicPath: process.env.NODE_ENV === 'production'
    ? ''
    : `http://localhost:3001/`,

  api: {
    uri: `http://localhost:4000/graphql`,
  },

  security: {
    secretKey: 'mysecretkey',
    expiresIn: '1d',
  },
};

const { Configuration } = require('webpack');

const webpackProdConfig = ({ presets }) => {
  const config = {
    mode: 'production',
  };
  return config;
};

module.exports = webpackProdConfig;

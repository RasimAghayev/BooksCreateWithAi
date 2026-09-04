const { getWebpackCommonConfig } = require('./webpack.common');
const { getWebpackDevelopmentConfig } = require('./webpack.development');
const { getWebpackProductionConfig } = require('./webpack.production');

const { merge } = require('webpack-merge');

const getWebpackConfig = (args) => {
  const configs = [getWebpackCommonConfig(args)];

  if (args.mode === 'development') {
    configs.push(getWebpackDevelopmentConfig());
  } else if (args.mode === 'production') {
    configs.push(getWebpackProductionConfig(args));
  }

  return merge(...configs);
};

module.exports = getWebpackConfig;

// CLI
const log = require('./cli/log');
// Webpack
const getWebpackConfig = require('./webpack/webpack.config');
const { getWebpackCommonConfig } = require('./webpack/webpack.common');
const { getWebpackDevelopmentConfig } = require('./webpack/webpack.development');
const { getWebpackProductionConfig } = require('./webpack/webpack.production');

module.exports = {
  log,
  getWebpackConfig,
  getWebpackCommonConfig,
  getWebpackDevelopmentConfig,
  getWebpackProductionConfig,
};

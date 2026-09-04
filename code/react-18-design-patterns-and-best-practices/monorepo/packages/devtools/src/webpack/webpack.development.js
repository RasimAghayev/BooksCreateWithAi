const webpack = require('webpack');

const getWebpackDevelopmentConfig = () => {
  const webpackConfig = {
    mode: 'development',
    devtool: 'source-map',
    plugins: [
      new webpack.HotModuleReplacementPlugin(),
      new webpack.NoEmitOnErrorsPlugin(),
    ],
  };
  return webpackConfig;
};

module.exports = { getWebpackDevelopmentConfig };

const { Configuration, HotModuleReplacementPlugin, NoEmitOnErrorsPlugin } = require('webpack');

const webpackDevConfig = () => {
  const config = {
    mode: 'development',
    devtool: 'source-map',
    output: {
      filename: '[name].js',
    },
    plugins: [
      new HotModuleReplacementPlugin(),
      new NoEmitOnErrorsPlugin(),
    ],
  };
  return config;
};

module.exports = webpackDevConfig;

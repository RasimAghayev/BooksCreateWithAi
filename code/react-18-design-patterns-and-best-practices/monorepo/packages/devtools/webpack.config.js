const webpack = require('webpack');
const nodeExternals = require('webpack-node-externals');
const path = require('path');

const createWebpackConfig = (env) => {
  const isDevelopment = env.mode === 'development';
  const isProduction = env.mode === 'production';

  const config = {
    mode: env.mode,
    resolve: {
      extensions: ['.ts', '.tsx', '.js', '.jsx', '.json'],
    },
    module: {
      rules: [
        {
          test: /\.[jt]sx?$/,
          exclude: /node_modules/,
          use: 'ts-loader',
        },
        {
          test: /\.css$/,
          use: ['style-loader', 'css-loader'],
        },
        {
          test: /\.svg$/,
          use: '@svgr/webpack',
        },
        {
          test: /\.(png|jpe?g|gif)$/i,
          type: 'asset/resource',
        },
      ],
    },
    plugins: [
      new webpack.ProgressPlugin(),
    ],
    externals: isProduction
      ? [nodeExternals()]
      : [],
  };

  if (isDevelopment) {
    config.devtool = 'source-map';
    config.devServer = {
      hot: true,
      port: 3000,
    };
  }

  return config;
};

module.exports = createWebpackConfig;

const Dotenv = require('dotenv-webpack');
const path = require('path');
const { Configuration } = require('webpack');

const styledComponentsTransformer = (() => ({
  before: [],
  after: [],
}))();

const webpackCommonConfig = () => {
  const config = {
    output: {
      path: path.resolve('dist'),
    },
    resolve: {
      extensions: ['.ts', '.tsx', '.js', '.jsx', '.json'],
      alias: {
        '~': path.resolve(__dirname, '../src'),
      },
    },
    module: {
      rules: [
        {
          test: /\.(woff|woff2)$/,
          type: 'asset/resource',
        },
        {
          test: /\.(ts|tsx)$/,
          exclude: /node_modules/,
          use: {
            loader: 'ts-loader',
            options: {
              transpileOnly: true,
            },
          },
        },
      ],
    },
    plugins: [new Dotenv()],
    optimization: {
      splitChunks: {
        cacheGroups: {
          default: false,
          vendor: {
            test: /node_modules/,
            name: 'vendor',
            chunks: 'all',
          },
        },
      },
    },
  };
  return config;
};

module.exports = webpackCommonConfig;

const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');
const WebpackBar = require('webpackbar');

const isAnalyze = Boolean(process.env.ANALYZE);

const webpackClientConfig = ({ mode }) => {
  const isProductionMode = mode === 'production';
  const config = {
    target: 'web',
    entry: {
      main: path.resolve(__dirname, '../../src/client/index.tsx'),
    },
    resolve: {
      extensions: ['.ts', '.tsx', '.js', '.jsx'],
      fallback: {
        buffer: require.resolve('buffer/'),
        crypto: require.resolve('crypto-browserify/'),
        stream: require.resolve('stream-browserify'),
      },
    },
    output: {
      path: path.resolve(__dirname, '../../dist/client'),
      filename: isProductionMode ? '[name].[contenthash].js' : '[name].js',
      clean: true,
    },
    module: {
      rules: [
        { test: /\.(ts|tsx)$/, exclude: /node_modules/, use: 'ts-loader' },
      ],
    },
    plugins: [
      new HtmlWebpackPlugin({
        template: path.resolve(__dirname, '../../src/client/index.html'),
      }),
      new WebpackBar(),
      ...(isAnalyze ? [new BundleAnalyzerPlugin()] : []),
    ],
    devServer: {
      port: 3000,
      historyApiFallback: true,
      hot: true,
    },
  };
  return config;
};

module.exports = webpackClientConfig;

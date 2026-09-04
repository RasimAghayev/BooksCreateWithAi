const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');

const isProduction = process.env.NODE_ENV === 'production';

module.exports = {
  // Generate source maps only in development
  devtool: !isProduction ? 'source-map' : false,
  entry: './src/index.tsx',
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: '[name].[hash:8].js',
    sourceMapFilename: '[name].[hash:8].map',
    chunkFilename: '[id].[hash:8].js',
    publicPath: '/',
  },
  resolve: {
    // Extensions supported as import paths
    extensions: ['.ts', '.tsx', '.js', '.json', '.css'],
  },
  target: 'web',
  mode: isProduction ? 'production' : 'development', // production minifies the bundle
  module: {
    rules: [
      {
        test: /\.(tsx|ts)$/,
        exclude: /node_modules/,
        use: {
          loader: 'ts-loader',
          options: {
            transpileOnly: true, // skip type-checking in webpack (handled by ForkTsChecker)
          },
        },
      },
      {
        test: /\.css/,
        use: [
          'style-loader', // injects CSS into <head>
          {
            loader: 'css-loader',
            options: {
              modules: {
                // dev-readable class names (hash in prod)
                localIdentName:
                  isProduction
                    ? '[hash:base64:5]'
                    : '[local]--[hash:base64:5]',
              },
            },
          },
        ],
      },
    ],
  },
  plugins: [
    new ForkTsCheckerWebpackPlugin(),
    new HtmlWebpackPlugin({
      title: 'Your project name',
      template: './src/index.html',
      filename: './index.html',
    }),
  ],
  optimization: {
    // Split vendor (node_modules) and main app bundles
    splitChunks: {
      cacheGroups: {
        default: false,
        commons: {
          test: /node_modules/,
          name: 'vendor',
          chunks: 'all',
        },
      },
    },
  },
};

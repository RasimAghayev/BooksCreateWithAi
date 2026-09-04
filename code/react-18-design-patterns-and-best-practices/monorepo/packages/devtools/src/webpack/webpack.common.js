const webpack = require('webpack');
const path = require('path');
const HtmlWebPackPlugin = require('html-webpack-plugin');
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');
const nodeExternals = require('webpack-node-externals');
const { log } = require('../cli/log');

const getWebpackCommonConfig = (args) => {
  const {
    configType,
    isAnalyze,
    port = 3000,
    mode,
    analyzerPort = 9001,
    packageName,
    htmlOptions,
    sandbox,
    devServer,
  } = args;

  const devServerPort = sandbox && devServer ? 8080 : port + 1;

  const entry = configType !== 'web'
    ? path.resolve(__dirname, `../../../${packageName}/src/index.ts`)
    : path.resolve(__dirname, `../../../${packageName}/src/index.tsx`);

  const resolve = {
    extensions: ['*', '.ts', '.tsx', '.js', '.jsx'],
    alias: {
      '~': path.resolve(__dirname, `../../../${packageName}/src`),
    },
    fallback: {
      buffer: false,
      crypto: false,
      stream: false,
      querystring: false,
      os: false,
      zlib: false,
      http: false,
      https: false,
      url: false,
      path: require.resolve('path-browserify'),
    },
  };

  const output = {
    path: path.resolve(__dirname, `../../../${packageName}/dist`),
    filename: '[name].js',
  };

  if (sandbox) {
    output.publicPath = '/';
    output.chunkFilename = '[name].js';
  }

  if (configType === 'package' && !sandbox) {
    output.filename = 'index.js';
    output.libraryTarget = 'umd';
    output.library = 'lib';
    output.umdNamedDefine = true;
    output.globalObject = 'this';
  }

  const plugins = [];

  if (isAnalyze) {
    plugins.push(new BundleAnalyzerPlugin({ analyzerPort }));
  }

  if (mode === 'development' && htmlOptions?.title && htmlOptions.template) {
    plugins.push(new HtmlWebPackPlugin({
      title: htmlOptions.title,
      template: path.resolve(__dirname, `../../../${packageName}/${htmlOptions.template}`),
      filename: './index.html',
    }));
  }

  const rules = [];

  rules.push({
    test: /\.(tsx|ts)$/,
    exclude: /node_modules/,
    loader: 'ts-loader',
    options: {
      transpileOnly: true,
    },
  });

  rules.push({
    test: /\.css$/,
    use: ['style-loader', 'css-loader'],
  });

  if (packageName === 'design-system') {
    rules.push({
      test: /\.svg$/,
      oneOf: [
        { use: 'svg-url-loader' },
        { use: '@svgr/webpack' },
      ],
    });
  }

  if (configType === 'package' && sandbox) {
    rules.push({
      test: /\.(jpe?g|png|gif|svg)$/i,
      use: [{ loader: 'file-loader' }],
    });
  }

  const webpackConfig = {
    entry,
    ...(sandbox && {
      entry: path.resolve(__dirname, `../../../${packageName}/sandbox/index.tsx`),
    }),
    ...(devServer && {
      devServer: {
        historyApiFallback: true,
        static: output.path,
        port: devServerPort,
      },
    }),
    ...(!sandbox && {
      externals: [nodeExternals()],
    }),
    output,
    resolve,
    plugins,
    module: { rules },
    ...(configType !== 'web' && !sandbox && {
      target: 'node',
    }),
  };

  log({ tag: 'webpack-common', json: webpackConfig, type: 'info' });
  return webpackConfig;
};

module.exports = { getWebpackCommonConfig };

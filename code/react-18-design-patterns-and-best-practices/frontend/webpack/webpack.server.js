const path = require('path');
const { Configuration, IgnorePlugin, optimize } = require('webpack');
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');
const nodeExternals = require('webpack-node-externals');
const WebpackBar = require('webpackbar');

const webpackServerConfig = ({ mode }) => {
  const isDevelopment = mode === 'development';
  const config = {
    target: 'node',
    entry: './src/server/index.ts',
    output: {
      libraryTarget: 'commonjs2',
      filename: 'server.js',
      path: path.resolve('dist'),
    },
    externals: [nodeExternals()],
    plugins: [
      new optimize.LimitChunkCountPlugin({ maxChunks: 1 }),
      new IgnorePlugin({
        resourceRegExp: /\.(sc|c)ss|jpe?g|png|gif|svg)$/i,
      }),
      new WebpackBar({
        name: 'server',
        color: '#2EA1F8',
        profile: true,
        basic: false,
      }),
    ],
  };

  if (isDevelopment) {
    config.watch = true;
    config.externals = [
      nodeExternals({
        allowlist: ['webpack/hot/poll?300'],
      }),
    ];
  }

  if (process.env.ANALYZE) {
    config.plugins.push(...(config.plugins || []), new BundleAnalyzerPlugin({ analyzerPort: 9002 }));
  }

  return config;
};

module.exports = webpackServerConfig;

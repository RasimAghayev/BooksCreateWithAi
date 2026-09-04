module.exports = {
  WebpackMode: 'production',
  ConfigType: 'web',
  Package: 'api',
  ConfigArgs: {
    mode: 'WebpackMode',
    type: 'ConfigType',
    packageName: 'Package',
  },
  ModeArgs: {
    configType: 'ConfigType',
    packageName: 'Package',
    mode: 'WebpackMode',
    devServer: false,
    isAnalyze: false,
    port: 3000,
    analyzerPort: 9001,
  },
};

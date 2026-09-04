const { ConfigArgs, getWebpackCommonConfig, getWebpackDevelopmentConfig, getWebpackProductionConfig, log } = require('@web-creator/devtools');
const { merge } = require('webpack-merge');

const getModeConfig = {
  development: getWebpackDevelopmentConfig,
  production: getWebpackProductionConfig,
};

const modeConfig = ({ mode, type, packageName }) => {
  const getWebpackConfiguration = getModeConfig[mode];
  return getWebpackConfiguration({
    configType: type,
    packageName,
    sandbox: true,
    devServer: true,
  });
};

const webpackConfig = async ({
  mode = 'production',
  type = 'web',
  sandbox = 'false',
  packageName = 'design-system',
} = {}) => {
  const isSandbox = type === 'package' && sandbox === 'true';
  const commonConfiguration = getWebpackCommonConfig({
    configType: type,
    packageName,
    mode,
    ...(isSandbox && {
      htmlOptions: { title: 'Sandbox', template: 'sandbox/index.html' },
      sandbox: isSandbox,
      devServer: isSandbox,
    }),
  });
  const modeConfiguration = mode && type ? modeConfig({ mode, type, packageName }) : {};
  const webpackConfiguration = merge(commonConfiguration, modeConfiguration);
  log({ tag: 'Webpack Configuration', json: webpackConfiguration, type: 'warning' });
  return webpackConfiguration;
};

module.exports = webpackConfig;

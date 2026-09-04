const { merge } = require('webpack-merge');

const loadPresets = async (env) => {
  const presets = [].concat(...(env.presets || []));
  const configs = await Promise.all(
    presets.map(async (presetName) => {
      try {
        const preset = require(`./presets/webpack.${presetName}`);
        return preset(env);
      } catch (err) {
        return {};
      }
    })
  );
  return merge({}, ...configs);
};

module.exports = loadPresets;

const cliColor = require('cli-color');

const log = (args) => {
  const blockColor = {
    info: cliColor.bgCyan.whiteBright.bold,
    error: cliColor.bgRed.whiteBright.bold,
    warning: cliColor.bgYellow.blackBright.bold,
  };
  const textColor = {
    info: cliColor.blue,
    error: cliColor.red,
    warning: cliColor.yellow,
  };

  if (typeof args === 'string') {
    console.info(textColor.info(args));
  }

  const { tag, json, type = 'info' } = args;
  if (tag && json) {
    console.info(blockColor[type](`<<< BEGIN ${tag.toUpperCase()}`));
    console.info(textColor[type](JSON.stringify(json, null, 2)));
    console.info(blockColor[type](`END ${tag.toUpperCase()} >>>`));
  }
};

module.exports = { log };

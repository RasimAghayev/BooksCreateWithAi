const { Buffer } = require('buffer');

global.Buffer = global.Buffer || Buffer;
global.window = global.window || {};

module.exports = { Buffer };

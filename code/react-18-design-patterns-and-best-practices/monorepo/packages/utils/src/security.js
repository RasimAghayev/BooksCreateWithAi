const bcrypt = require('bcryptjs');

const encrypt = (password) => {
  return bcrypt.hashSync(password, 10);
};

const decrypt = (password, hash) => {
  return bcrypt.compareSync(password, hash);
};

module.exports = { encrypt, decrypt };

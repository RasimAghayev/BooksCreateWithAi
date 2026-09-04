const { getUserBy, getUserData } = require('../../lib/auth');

const getUsers = (_, _args, { models }) => {
  return models.User.findAll();
};

const getUser = async (_, { at }, { models }) => {
  const connectedUser = await getUserData(at);
  if (connectedUser) {
    const user = await getUserBy(
      {
        id: connectedUser.id,
        email: connectedUser.email,
        active: connectedUser.active,
      },
      [connectedUser.role],
      models
    );
    if (user) {
      return connectedUser;
    }
  }
  return {
    id: '',
    username: '',
    email: '',
    role: '',
    active: false,
  };
};

const createUser = (_, { input }, { models }) => {
  return models.User.create({ ...input });
};

const login = (_, { input }, { models }) => {
  return authenticate(input.emailOrUsername, input.password, models);
};

const authenticate = async (emailOrUsername, password, models) => {
  const auth = require('../../lib/auth');
  return auth.authenticate(emailOrUsername, password, models);
};

module.exports = {
  Query: {
    getUser,
    getUsers,
  },
  Mutation: {
    createUser,
    login,
  },
};

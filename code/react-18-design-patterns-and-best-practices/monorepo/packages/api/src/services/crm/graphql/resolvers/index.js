const { mergeResolvers } = require('@graphql-tools/merge');
const userResolver = require('../../../../graphql/resolvers/user');
const guestResolver = require('./guest');

const resolvers = mergeResolvers([userResolver, guestResolver]);

module.exports = resolvers;

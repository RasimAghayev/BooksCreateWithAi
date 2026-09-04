const { mergeTypeDefs } = require('@graphql-tools/merge');
const ErrorType = require('../../../../graphql/types/Error');
const ScalarType = require('../../../../graphql/types/Scalar');
const UserType = require('../../../../graphql/types/User');
const GuestType = require('./Guest');

const typeDefs = mergeTypeDefs([ErrorType, ScalarType, UserType, GuestType]);

module.exports = typeDefs;

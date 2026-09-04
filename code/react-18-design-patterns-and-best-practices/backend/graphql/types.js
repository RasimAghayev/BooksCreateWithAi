// GraphQL schema: types, queries, mutations, input types, type-merging
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 13 (pp. 296-298)

import { gql } from 'graphql-tag';
import { mergeTypeDefs } from '@graphql-tools/merge';

// Scalar (custom) types
const Scalar = gql`
  scalar UUID
  scalar Datetime
  scalar JSON
`;

// Object type (note: GraphQL types are capitalized; ! = non-nullable)
const User = gql`
  type User {
    id: UUID!
    username: String!
    email: String!
    password: String!
    role: String!
    @active: Boolean!
    createdAt: Datetime!
    updatedAt: Datetime!
  }
`;

// Query type (read/fetch)
const Query = gql`
  type Query {
    getUser(at: String!): User!
    getUsers: [User!]
  }
`;

// Mutation type (write) + input types
const Mutation = gql`
  type Mutation {
    createUser(input: CreateUserInput): User!
    login(input: LoginInput): Token!
  }

  type Token {
    token: String!
  }

  input CreateUserInput {
    username: String!
    password: String!
    email: String!
    active: Boolean!
    role: String!
  }

  input LoginInput {
    emailOrUsername: String!
    password: String!
  }
`;

// Merge all type-defs into one schema
const typeDefs = mergeTypeDefs([Scalar, User, Query, Mutation]);
export default typeDefs;

/*
Resolvers receive (_parent, args, context, info):
  Query: {
    getUsers: (_, args, { models }) => models.User.findAll(),
    getUser: (_, { at }, { models }) => models.User.findOne(...)
  },
  Mutation: {
    createUser: (_, { input }, { models }) => models.User.create({ ...input }),
    login: (_, { input }, { models }) => doLogin(input.email, input.password, models)
  }
*/

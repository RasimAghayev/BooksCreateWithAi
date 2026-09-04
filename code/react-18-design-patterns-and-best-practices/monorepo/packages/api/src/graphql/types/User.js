const gql = require('graphql-tag');

const UserType = gql`
  "User type"
  type User {
    id: UUID!
    username: String!
    email: String!
    role: String!
    active: Boolean!
    createdAt: Datetime!
    updatedAt: Datetime!
  }

  "Token type"
  type Token {
    token: String!
  }

  "User Query"
  type Query {
    getUser(at: String!): User!
    getUsers: [User!]
  }

  "User Mutation"
  type Mutation {
    createUser(input: ICreateUser): User!
    login(input: ILogin): Token!
  }

  "CreateUser Input"
  input ICreateUser {
    username: String!
    password: String!
    email: String!
    active: Boolean!
    role: String!
  }

  "Login Input"
  input ILogin {
    emailOrUsername: String!
    password: String!
  }
`;

module.exports = UserType;

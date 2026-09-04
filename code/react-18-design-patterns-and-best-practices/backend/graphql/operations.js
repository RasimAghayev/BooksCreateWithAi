const { gql } = require('@apollo/client');

const GET_USERS = gql`
  query {
    getUsers {
      id
      username
      email
      role
    }
  }
`;

const GET_USER = gql`
  query {
    getUser {
      id
      username
      email
      role
    }
  }
`;

const CREATE_USER = gql`
  mutation($input: CreateUserInput) {
    createUser(input: $input) {
      id
      username
      email
      role
      active
    }
  }
`;

const LOGIN_USER = gql`
  mutation($input: LoginInput) {
    login(input: $input) {
      token
      user {
        id
        username
        email
        role
      }
    }
  }
`;

module.exports = { GET_USERS, GET_USER, CREATE_USER, LOGIN_USER };

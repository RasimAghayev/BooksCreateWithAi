const { gql } = require('@apollo/client');

const LOGIN_MUTATION = gql`
  mutation login($email: String!, $password: String!) {
    login(input: { email: $email, password: $password }) {
      token
    }
  }
`;

module.exports = LOGIN_MUTATION;

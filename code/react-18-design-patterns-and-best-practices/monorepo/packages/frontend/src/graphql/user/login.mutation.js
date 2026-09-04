const gql = require('@apollo/client').gql;

const LoginMutation = gql`
  mutation login($emailOrUsername: String!, $password: String!) {
    login(input: { emailOrUsername: $emailOrUsername, password: $password }) {
      token
    }
  }
`;

module.exports = LoginMutation;

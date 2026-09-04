const { gql } = require('@apollo/client');

const GET_USER_QUERY = gql`
  query getUser($at: String!) {
    getUser(at: $at) {
      id
      email
      username
      role
      active
    }
  }
`;

module.exports = GET_USER_QUERY;

const gql = require('@apollo/client').gql;

const GetUserQuery = gql`
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

module.exports = GetUserQuery;

const gql = require('graphql-tag');

const ErrorType = gql`
  type ErrorResponse {
    code: Int
    message: String!
  }

  type Error {
    error: ErrorResponse
  }
`;

module.exports = ErrorType;

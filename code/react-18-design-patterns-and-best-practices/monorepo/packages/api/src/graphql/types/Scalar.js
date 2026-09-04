const gql = require('graphql-tag');

const ScalarType = gql`
  scalar UUID
  scalar Datetime
  scalar JSON
`;

module.exports = ScalarType;

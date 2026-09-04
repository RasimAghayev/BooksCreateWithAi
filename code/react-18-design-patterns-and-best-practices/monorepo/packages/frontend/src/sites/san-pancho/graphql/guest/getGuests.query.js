const gql = require('@apollo/client').gql;

const getGuestsQuery = `
  getGuests {
    on GuestResponse {
      guests {
        id
        fullName
        email
        photo
        socialMedia
        location
        gender
        birthday
      }
    }
    on Error {
      error {
        code
        message
      }
    }
  }
`;

const query = gql`
  ${getGuestsQuery}
`;

const getGuests = gql`
  query getGuests {
    ${getGuestsQuery}
  }
`;

module.exports = { getGuests, getGuestsQuery };

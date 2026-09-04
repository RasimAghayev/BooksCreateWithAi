const guestResolver = {
  Query: {
    getGuests: async (_, _args, { models }) => {
      const guests = await models.Guest.findAll({
        order: [['fullName', 'ASC']],
      });
      if (guests.length > 0) {
        return {
          __typename: 'GuestResponse',
          guests,
        };
      }
      return {
        __typename: 'Error',
        error: {
          code: 404,
          message: 'No guests found',
        },
      };
    },
  },
};

module.exports = guestResolver;

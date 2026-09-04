// Sequelize User model (PostgreSQL) with validation + beforeCreate password encryption
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 13 (pp. 303-306)

import { encrypt } from '@contentpi/lib';

// IDataTypes mirrors Sequelize's DataTypes; passed by models/index.ts
// sequelize.define('User', { fields }, { hooks })
export default (sequelize, DataTypes) => {
  const User = sequelize.define(
    'User',
    {
      id: {
        primaryKey: true,
        allowNull: false,
        type: DataTypes.UUID,
        defaultValue: DataTypes.UUIDV4(),
      },
      username: {
        type: DataTypes.STRING,
        allowNull: false,
        unique: true,
        validate: {
          isAlphanumeric: {
            args: true,
            msg: 'The user just accepts alphanumeric characters',
          },
          len: {
            args: [4, 20],
            msg: 'The username must be from 4 to 20 characters',
          },
        },
      },
      password: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      email: {
        type: DataTypes.STRING,
        allowNull: false,
        unique: true,
        validate: {
          isEmail: { args: true, msg: 'Invalid email' },
        },
      },
      role: {
        type: DataTypes.STRING,
        allowNull: false,
        defaultValue: 'user',
      },
      active: {
        type: DataTypes.BOOLEAN,
        allowNull: false,
        defaultValue: false,
      },
    },
    {
      hooks: {
        // Encrypt password (sha1) BEFORE INSERT
        beforeCreate: (user) => {
          user.password = encrypt(user.password);
        },
      },
    }
  );
  return User;
};

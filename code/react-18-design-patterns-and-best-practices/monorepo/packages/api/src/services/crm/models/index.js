require('module-alias/register');

const pg = require('pg');
const { Sequelize } = require('sequelize');
const { keys, ts } = require('@web-creator/utils');
const Config = require('../../../config').Config;
const createUserModel = require('../../../models/User');
const createGuestModel = require('./Guest');

const { engine, port, host, database, username, password } = Config.database || {};
const uri = `${engine}://${username}:${password}@${host}:${port}/${database}`;
const sequelize = new Sequelize(uri, {
  dialectModule: pg,
});

const addModel = (path) => require(path).default(sequelize, Sequelize);

const models = {
  User: createUserModel(sequelize, Sequelize),
  Guest: createGuestModel(sequelize, Sequelize),
  sequelize,
};

const keys = Object.keys(models);
keys.forEach((modelName) => {
  if (models[modelName].associate) {
    models[modelName].associate(models);
  }
});

module.exports = models;

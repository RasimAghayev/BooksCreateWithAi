const pg = require('pg');
const { Sequelize } = require('sequelize');
const Config = require('../../../config').Config;
const createUserModel = require('../../../../models/User');
const createGuestModel = require('./Guest');

const { engine, port, host, database, username, password } = Config.database || {};
const uri = `${engine}://${username}:${password}@${host}:${port}/${database}`;
const sequelize = new Sequelize(uri, {
  dialectModule: pg,
});

const models = {
  User: createUserModel(sequelize, Sequelize),
  sequelize,
};

module.exports = models;

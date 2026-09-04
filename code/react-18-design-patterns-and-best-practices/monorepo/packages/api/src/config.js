const dotenv = require('dotenv');
const { config: crmConfig } = require('./services/crm/config');
const { config: blankServiceConfig } = require('./services/default/config');

dotenv.config();

const getServiceConfig = (service) => {
  switch (service) {
    case 'crm':
      return crmConfig;
    default:
      return blankServiceConfig;
  }
};

const buildConfig = () => {
  const service = process.env.SERVICE;
  if (!service) {
    throw 'You must specify a service (E.g., SERVICE=crm npm run dev)';
  }
  const serviceConfig = getServiceConfig(service);
  const config = {
    ...serviceConfig,
    service,
    database: {
      ...serviceConfig.database,
      engine: process.env.DB_ENGINE,
      host: process.env.DB_HOST,
      port: Number(process.env.DB_PORT),
      username: process.env.DB_USERNAME,
      password: process.env.DB_PASSWORD,
    },
  };
  return config;
};

const Config = buildConfig();
module.exports = { Config, getServiceConfig, buildConfig };

const { makeExecutableSchema } = require('@graphql-tools/schema');
const { ApolloServer } = require('@apollo/server');
const { expressMiddleware } = require('@apollo/server/express4');
const { ApolloServerPluginDrainHttpServer } = require('@apollo/server/plugin/drainHttpServer');
const bodyParser = require('body-parser');
const http = require('http');
const cookieParser = require('cookie-parser');
const cors = require('cors');
const express = require('express');
const { applyMiddleware } = require('graphql-middleware');
const { json } = require('body-parser');
const { ts } = require('@web-creator/utils');
const { Service } = require('./types/config');

const service = process.env.SERVICE || 'default';

const serviceNames = ['CRM', 'default'];
if (!ts.includes(serviceNames, service)) {
  throw 'Invalid service';
}

const resolvers = require(`./services/${service}/graphql/resolvers`).default || require(`./services/${service}/graphql/resolvers`);
const typeDefs = require(`./services/${service}/graphql/types`).default || require(`./services/${service}/graphql/types`);
const models = require(`./services/${service}/models`).default || require(`./services/${service}/models`);
const seeds = require(`./services/${service}/seeds`).default || require(`./services/${service}/seeds`);

const app = express();
const httpServer = http.createServer(app);

const corsOptions = {
  origin: '*',
  credentials: true,
};

app.use(cors(corsOptions));
app.use(cookieParser());
app.use(bodyParser.json());
app.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Headers', 'Origin, X-Requested-With, Content-Type, Accept');
  next();
});

const schema = applyMiddleware(
  makeExecutableSchema({
    typeDefs,
    resolvers,
  })
);

const apolloServer = new ApolloServer({
  schema,
  plugins: [ApolloServerPluginDrainHttpServer({ httpServer })],
});

const main = async () => {
  const alter = true;
  const force = false;

  await apolloServer.start();
  await models.sequelize.sync({ alter, force });
  console.log('Initializing Seeds...');
  seeds();

  app.use(
    '/graphql',
    cors(),
    json(),
    expressMiddleware(apolloServer, {
      context: async () => ({ models }),
    })
  );

  await new Promise((resolve) => httpServer.listen({ port: 4000 }, resolve));
  console.log('Server ready at http://localhost:4000/graphql');
};

main();

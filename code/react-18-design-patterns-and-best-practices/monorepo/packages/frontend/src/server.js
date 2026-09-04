const express = require('express');
const nextJS = require('next');
const path = require('path');
const cookieParser = require('cookie-parser');
const { ts } = require('@web-creator/utils');
const Config = require('./config').Config;
const { Site } = require('./types/config');

const site = process.env.SITE || 'blank-page';

const siteNames = ['san-pancho', 'codejobs', 'blank-page'];
if (!ts.includes(siteNames, site)) {
  throw 'Invalid site';
}

const { hostname } = Config;
const port = 3000;
const dev = process.env.NODE_ENV !== 'production';
const nextApp = nextJS({ dev, hostname, port });
const handle = nextApp.getRequestHandler();

nextApp.prepare().then(() => {
  const app = express();

  app.use(cookieParser());
  app.use(express.static(path.join(__dirname, '../public')));
  app.use(express.static(path.join(__dirname, `./sites/${Config.site}/static`)));

  app.get('/logout', (req, res) => {
    const redirect = req.query.redirectTo || '/';
    res.clearCookie(`at-${Config.site}`);
    res.redirect(redirect);
  });

  app.get('/dashboard', (req, res, next) => {
    if (Config.user) {
      next();
    } else {
      res.redirect(`/login?redirectTo=/dashboard`);
    }
  });

  app.all('*', (req, res) => handle(req, res));
  app.listen(3000);
});

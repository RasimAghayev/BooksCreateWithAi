const cookieParser = require('cookie-parser');
const cors = require('cors');
const express = require('express');
const path = require('path');
const config = require('../config');
const { isConnected } = require('./lib/middlewares/user');

const app = express();
const distDir = path.resolve('dist');
const staticDir = path.resolve('src', 'static');

app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(cookieParser(config.security.secretKey));
app.use(cors({ credentials: true, origin: true }));

app.use(express.static(distDir));
app.use(express.static(staticDir));

// Protected route: only god/admin users
app.get('/dashboard', isConnected(true), (req, res) => {
  res.send('dashboard');
});

// Login route: redirect if already connected
app.get('/login', isConnected(false), (req, res) => {
  res.send('login');
});

// Logout: clear session cookie
app.get('/logout', (req, res) => {
  const redirect = req.query.redirectTo || '/';
  res.clearCookie('at');
  res.redirect(redirect);
});

// Catch-all: let React handle routing
app.get('*', (req, res) => {
  res.send('app');
});

module.exports = app;

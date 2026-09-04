module.exports = {
  apps: [{
    name: 'my-app',
    script: 'node',
    args: 'server.js',
    cwd: '/var/www/my-app',
    instances: 'max',
    exec_mode: 'cluster',
    env: {
      NODE_ENV: 'development',
      PORT: 3000,
    },
    env_production: {
      NODE_ENV: 'production',
      PORT: 3000,
    },
  }],
};

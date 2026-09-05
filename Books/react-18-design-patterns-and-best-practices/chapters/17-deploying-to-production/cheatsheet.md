# Cheat Sheet — Deploying to Production

## Deployment

### nginx reverse proxy

**Nə edir:** nginx Node.js-ə reverse proxy kimi qoşulur.

**Kod:**

```nginx
server {
  listen 80;
  server_name example.com;
  location / {
    proxy_pass http://localhost:3000;
    proxy_set_header Host $host;
  }
}
```

**Mənbə:** Chapter 17, page 467-500

### PM2 start

**Nə edir:** Node.js-i background-da işlədir. `pm2 save` və `pm2 startup` ilə reboot davamlılığı.

**Kod:**

```shell
pm2 start npm --name 'app' -- start
```

**Mənbə:** Chapter 17, page 467-500

### Let's Encrypt SSL

**Nə edir:** Pulsuz SSL sertifikatı. Nginx konfiqurasiyasını avtomatik yeniləyir.

**Kod:**

```shell
sudo certbot --nginx -d example.com
```

**Mənbə:** Chapter 17, page 467-500

### CircleCI config.yml

**Nə edir:** CI/CD pipeline — lint, test, build, deploy.

**Kod:**

```yaml
version: 2.1
jobs:
  build:
    docker: [{ image: cimg/node:18.0 }]
    steps:
      - checkout
      - run: npm ci
      - run: npm test
      - run: npm run build
      - run: ssh user@$DROPLET_IP 'pm2 restart app'
```

**Mənbə:** Chapter 17, page 467-500

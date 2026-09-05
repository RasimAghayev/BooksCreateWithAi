# 17. Deploying to Production

**Səhifələr:** 467-500

## Bu fəsil nədən bəhs edir?

Bu fəsil React tətbiqini production mühitə necə yerləşdirməyi öyrədir: DigitalOcean Droplet, nginx, PM2, domain qurma və CircleCI ilə CI/CD.

## Əsas fikirlər

### 1. DigitalOcean Droplet yaratmaq

DigitalOcean-da Ubuntu Droplet yarat. SSH açarı əlavə et, root girişi konfiqurasiya et. Domain name qeydiyyatdan keçir (Namecheap, GoDaddy) və Droplet IP-yə yönləndir.

### 2. Node.js qurma

Server-ə Node.js-in LTS versiyasını qur (NVM və ya system package). Layihəni GitHub-dan clone et, `npm ci` ilə dependencies qur, `.env` faylını server-ə yüklə (secret management).

### 3. nginx qurma

nginx reverse proxy kimi işləyir: brauzerdən gələn sorğuları Node.js-ə yönləndirir (port 3000). HTTPS üçün Let's Encrypt ilə SSL sertifikat al. nginx konfiqurasiya faylı `/etc/nginx/sites-available/default`.

### 4. PM2 ilə process idarəsi

PM2 Node.js prosesini background-da işlədir, crash olarsa avtomatik restart edir. `pm2 start npm --name 'app' -- start`, `pm2 save`, `pm2 startup` (system reboot-da avtomatik başlama). Log-lar `pm2 logs` ilə izlənir.

### 5. CircleCI ilə CI/CD

`/.circleci/config.yml` faylı yazılır. Pipeline mərhələri: install dependencies → lint → test → build → deploy. SSH açarı CircleCI-ə əlavə olunur ki, server-ə deploy edə bilsin. Environment variables CircleCI dashboard-da saxlanılır (DROPLET_IP, SSH_PRIVATE_KEY).

### 6. Production checklist

HTTPS aktivdir, gzip/compression aktiv, environment variables konfiqurasiya edilib, error tracking (Sentry) qurulub, backup strategiysı mövcuddur, monitoring (PM2 Plus, Datadog) qoşulub.

## Əsas terminlər

- Production
- Deployment
- DigitalOcean
- Droplet
- Ubuntu
- Node.js
- nginx
- PM2
- Reverse Proxy
- HTTPS
- Let's Encrypt
- CircleCI
- CI/CD
- SSH
- Domain

## Praktik nəticə

Kiçik bir React tətbiqini DigitalOcean Droplet-ə yerləşdir: nginx qur, PM2 ilə start et, CircleCI ilə avtomatik deploy qur.

## Mənbə

Pages: 467-500
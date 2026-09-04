# Deployment Checklist

## 1. DigitalOcean Droplet
- [ ] Create Droplet (Ubuntu 20.04, $6/mo, Basic plan)
- [ ] SSH into Droplet: `ssh root@YOUR_DROPLET_IP`

## 2. Node.js + PM2
- [ ] Install Node.js: `curl -sL https://deb.nodesource.com/setup_19.x | sudo bash - && sudo apt install nodejs -y`
- [ ] Install PM2: `sudo npm install -g pm2`

## 3. Git + SSH
- [ ] Generate SSH key: `ssh-keygen -t rsa`
- [ ] Add public key to GitHub account
- [ ] Clone repo: `git clone git@github.com:...`

## 4. nginx
- [ ] Install: `sudo apt-get install nginx`
- [ ] Configure firewall: `sudo ufw allow 'Nginx HTTP'`
- [ ] Edit `/etc/nginx/sites-available/default` → set proxy_pass to localhost:3000
- [ ] Test: `sudo nginx -t`
- [ ] Restart: `sudo systemctl restart nginx`

## 5. Domain (optional)
- [ ] Point nameservers to DigitalOcean DNS
- [ ] Add domain in DO networking
- [ ] Create CNAME record (www → @)

## 6. PM2
- [ ] Start app: `pm2 start ecosystem.config.js --env production`
- [ ] Save: `pm2 save`
- [ ] Generate startup: `pm2 startup`

## 7. CircleCI (optional)
- [ ] Sign up at circleci.com with GitHub
- [ ] Add SSH key to CircleCI project settings
- [ ] Add env vars: `DROPLET_USER`, `DROPLET_IP`
- [ ] Commit `.circleci/config.yml`

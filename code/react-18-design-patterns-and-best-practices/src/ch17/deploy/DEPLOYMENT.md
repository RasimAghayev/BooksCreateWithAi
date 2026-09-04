# Deployment Summary

## Prerequisites
- DigitalOcean account
- Domain name (GoDaddy, Namecheap, etc.)

## Infrastructure
- [x] DigitalOcean Droplet (Ubuntu 20.04, $6/mo)
- [x] Node.js 19.x installed via nodesource PPA
- [x] PM2 installed globally for process management
- [x] nginx installed for reverse proxy
- [x] UFW firewall configured (port 80 open)
- [x] Domain nameservers pointed to DigitalOcean DNS
- [x] CNAME record created (www → @)

## CI/CD Pipeline (CircleCI)
1. Checkout code from GitHub
2. `npm install` — install dependencies
3. `npm run lint` — ESLint validation
4. `npm test` — Jest/Vitest tests
5. SSH deploy to Droplet:
   - `git checkout master`
   - `git pull`
   - `npm install`
   - `npm run start:production`

## Environment Variables (CircleCI)
- `DROPLET_USER` = root
- `DROPLET_IP` = your-droplet-ip

## Useful Commands
```bash
# On Droplet
ssh root@YOUR_DROPLET_IP
pm2 list                    # View running processes
pm2 logs                    # View logs
pm2 restart all            # Restart all processes
pm2 stop all               # Stop all processes

# nginx
sudo nginx -t              # Test config
sudo systemctl restart nginx  # Restart nginx
sudo ufw allow 'Nginx HTTP'   # Open port 80

# Local dev
npm start                  # Dev mode (port 3000)
npm run start:production   # Production via PM2
```

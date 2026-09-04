#!/bin/bash
# Setup Node.js on Ubuntu Droplet
set -e

echo "Setting up Node.js..."

cd ~
curl -sL https://deb.nodesource.com/setup_19.x -o nodesource_setup.sh
sudo bash nodesource_setup.sh
sudo apt install nodejs -y

echo "Node.js version:"
node -v
echo "npm version:"
npm -v

# Install PM2 globally
sudo npm install -g pm2

echo "Setup complete!"

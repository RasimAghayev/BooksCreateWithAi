#!/usr/bin/env bash
set -euo pipefail

npm install -g @babel/core @babel/node
babel source.js -o output.js
npm install -g @babel/preset-env @babel/preset-react

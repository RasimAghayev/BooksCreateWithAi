#!/usr/bin/env bash
# Chapter 1 — Vite setup (page 40)
set -euo pipefail

# React is split into two packages:
#   react       -> core features (shared across DOM / Native)
#   react-dom   -> browser-related features
# UMD fallbacks (no package manager):
#   https://unpkg.com/react@18.2.0/umd/react.production.min.js
#   https://unpkg.com/react-dom@18.2.0/umd/react-dom.production.min.js

npm install -g create-vite
create-vite my-react-app --template react-ts
cd my-react-app
npm install
npm run dev   # runs on port 5173 by default

#!/usr/bin/env node

const { execFileSync } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

const binName = os.platform() === 'win32' ? 'cl.exe' : 'cl';
const binPath = path.join(__dirname, binName);

if (!fs.existsSync(binPath)) {
  console.error('cl binary not found. Please reinstall the package.');
  console.error('Run: npm install -g @groo.dev/cl-wrangler');
  process.exit(1);
}

try {
  execFileSync(binPath, process.argv.slice(2), { stdio: 'inherit' });
} catch (err) {
  process.exit(err.status || 1);
}

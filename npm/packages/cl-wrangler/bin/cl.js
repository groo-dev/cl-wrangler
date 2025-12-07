#!/usr/bin/env node

const { execFileSync } = require('child_process');
const path = require('path');
const os = require('os');

const PLATFORMS = {
  'darwin-arm64': '@groo.dev/cl-wrangler-darwin-arm64',
  'darwin-x64': '@groo.dev/cl-wrangler-darwin-x64',
  'linux-arm64': '@groo.dev/cl-wrangler-linux-arm64',
  'linux-x64': '@groo.dev/cl-wrangler-linux-x64',
  'win32-arm64': '@groo.dev/cl-wrangler-win32-arm64',
  'win32-x64': '@groo.dev/cl-wrangler-win32-x64',
};

function getBinaryPath() {
  const platform = `${os.platform()}-${os.arch()}`;
  const pkg = PLATFORMS[platform];

  if (!pkg) {
    throw new Error(`Unsupported platform: ${platform}`);
  }

  const binName = os.platform() === 'win32' ? 'cl.exe' : 'cl';

  // Try to find the platform package
  try {
    const pkgPath = require.resolve(`${pkg}/package.json`);
    return path.join(path.dirname(pkgPath), 'bin', binName);
  } catch (e) {
    throw new Error(
      `Platform package ${pkg} not found. ` +
      `Please reinstall @groo.dev/cl-wrangler or install the platform package manually.`
    );
  }
}

try {
  const binPath = getBinaryPath();
  execFileSync(binPath, process.argv.slice(2), { stdio: 'inherit' });
} catch (err) {
  if (err.status !== undefined) {
    process.exit(err.status);
  }
  console.error(err.message);
  process.exit(1);
}

const fs = require('fs');
const os = require('os');
const path = require('path');
const https = require('https');

const PLATFORMS = {
  'darwin-arm64': '@groo.dev/cl-wrangler-darwin-arm64',
  'darwin-x64': '@groo.dev/cl-wrangler-darwin-x64',
  'linux-arm64': '@groo.dev/cl-wrangler-linux-arm64',
  'linux-x64': '@groo.dev/cl-wrangler-linux-x64',
  'win32-arm64': '@groo.dev/cl-wrangler-win32-arm64',
  'win32-x64': '@groo.dev/cl-wrangler-win32-x64',
};

const ARCHIVES = {
  'darwin-arm64': 'cl_darwin_arm64.tar.gz',
  'darwin-x64': 'cl_darwin_amd64.tar.gz',
  'linux-arm64': 'cl_linux_arm64.tar.gz',
  'linux-x64': 'cl_linux_amd64.tar.gz',
  'win32-arm64': 'cl_windows_arm64.zip',
  'win32-x64': 'cl_windows_amd64.zip',
};

function getPlatformKey() {
  return `${os.platform()}-${os.arch()}`;
}

function getBinName() {
  return os.platform() === 'win32' ? 'cl.exe' : 'cl';
}

function getVersion() {
  const pkg = require('./package.json');
  return pkg.version;
}

function getBinaryPathFromPackage() {
  const platform = getPlatformKey();
  const pkg = PLATFORMS[platform];

  if (!pkg) {
    throw new Error(`Unsupported platform: ${platform}`);
  }

  try {
    const pkgPath = require.resolve(`${pkg}/package.json`);
    return path.join(path.dirname(pkgPath), 'bin', getBinName());
  } catch (e) {
    return null;
  }
}

function downloadBinary(url) {
  return new Promise((resolve, reject) => {
    const follow = (url) => {
      https.get(url, (res) => {
        if (res.statusCode === 301 || res.statusCode === 302) {
          return follow(res.headers.location);
        }
        if (res.statusCode !== 200) {
          return reject(new Error(`Download failed: ${res.statusCode}`));
        }
        const chunks = [];
        res.on('data', (chunk) => chunks.push(chunk));
        res.on('end', () => resolve(Buffer.concat(chunks)));
        res.on('error', reject);
      }).on('error', reject);
    };
    follow(url);
  });
}

async function extractTarGz(buffer, destDir) {
  const zlib = require('zlib');
  const { execSync } = require('child_process');

  const tempFile = path.join(os.tmpdir(), `cl-${Date.now()}.tar.gz`);
  fs.writeFileSync(tempFile, buffer);

  try {
    execSync(`tar -xzf "${tempFile}" -C "${destDir}"`, { stdio: 'pipe' });
  } finally {
    fs.unlinkSync(tempFile);
  }
}

async function extractZip(buffer, destDir) {
  const { execSync } = require('child_process');

  const tempFile = path.join(os.tmpdir(), `cl-${Date.now()}.zip`);
  fs.writeFileSync(tempFile, buffer);

  try {
    if (os.platform() === 'win32') {
      execSync(`powershell -command "Expand-Archive -Path '${tempFile}' -DestinationPath '${destDir}' -Force"`, { stdio: 'pipe' });
    } else {
      execSync(`unzip -o "${tempFile}" -d "${destDir}"`, { stdio: 'pipe' });
    }
  } finally {
    fs.unlinkSync(tempFile);
  }
}

async function downloadAndInstall() {
  const platform = getPlatformKey();
  const archive = ARCHIVES[platform];
  const version = getVersion();

  if (!archive) {
    throw new Error(`Unsupported platform: ${platform}`);
  }

  const url = `https://github.com/groo-dev/cl-wrangler/releases/download/v${version}/${archive}`;
  console.error(`[cl-wrangler] Downloading binary from ${url}`);

  const buffer = await downloadBinary(url);
  const binDir = path.join(__dirname, 'bin');

  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  if (archive.endsWith('.zip')) {
    await extractZip(buffer, binDir);
  } else {
    await extractTarGz(buffer, binDir);
  }

  const binPath = path.join(binDir, getBinName());
  if (os.platform() !== 'win32') {
    fs.chmodSync(binPath, 0o755);
  }

  console.error(`[cl-wrangler] Binary installed successfully`);
  return binPath;
}

function optimizeBinary(binPath) {
  // Replace JS wrapper with hard link to binary for faster startup
  // Only on Unix (Windows doesn't support this well)
  if (os.platform() === 'win32') {
    return;
  }

  const jsWrapper = path.join(__dirname, 'bin', 'cl.js');
  const tempPath = path.join(__dirname, 'bin', 'cl-temp');

  try {
    // Create hard link from binary to temp location
    fs.linkSync(binPath, tempPath);
    // Replace the JS wrapper with the hard link
    fs.renameSync(tempPath, jsWrapper);
    console.error(`[cl-wrangler] Optimized: replaced JS wrapper with binary`);
  } catch (e) {
    // Optimization failed, but that's okay - JS wrapper still works
    try {
      fs.unlinkSync(tempPath);
    } catch {}
  }
}

async function main() {
  // Check if binary is available from platform package
  let binPath = getBinaryPathFromPackage();

  if (binPath && fs.existsSync(binPath)) {
    console.error(`[cl-wrangler] Found binary in platform package`);
  } else {
    // Fallback: download binary directly
    console.error(`[cl-wrangler] Platform package not found, downloading binary...`);
    console.error(`[cl-wrangler] This can happen if you used --no-optional or --ignore-optional`);
    binPath = await downloadAndInstall();
  }

  // Optimize by replacing JS wrapper with binary
  optimizeBinary(binPath);
}

main().catch((err) => {
  console.error(`[cl-wrangler] Installation failed: ${err.message}`);
  console.error(`[cl-wrangler] You can install manually from: https://github.com/groo-dev/cl-wrangler/releases`);
  process.exit(1);
});

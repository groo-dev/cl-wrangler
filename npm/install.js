const os = require('os');
const fs = require('fs');
const path = require('path');
const https = require('https');
const { execSync } = require('child_process');

const VERSION = process.env.CL_VERSION || '0.0.0';
const REPO = 'groo-dev/cl-wranger';

const PLATFORMS = {
  'darwin-x64': 'cl_darwin_amd64.tar.gz',
  'darwin-arm64': 'cl_darwin_arm64.tar.gz',
  'linux-x64': 'cl_linux_amd64.tar.gz',
  'linux-arm64': 'cl_linux_arm64.tar.gz',
  'win32-x64': 'cl_windows_amd64.zip',
  'win32-arm64': 'cl_windows_arm64.zip',
};

const platform = `${os.platform()}-${os.arch()}`;
const archive = PLATFORMS[platform];

if (!archive) {
  console.error(`Unsupported platform: ${platform}`);
  console.error('Please install manually from: https://github.com/groo-dev/cl-wranger/releases');
  process.exit(1);
}

const url = `https://github.com/${REPO}/releases/download/cli-v${VERSION}/${archive}`;
const binDir = path.join(__dirname, 'bin');
const binName = os.platform() === 'win32' ? 'cl.exe' : 'cl';
const binPath = path.join(binDir, binName);

// Skip if binary already exists (for development)
if (fs.existsSync(binPath)) {
  console.log('cl binary already exists, skipping download');
  process.exit(0);
}

console.log(`Downloading cl v${VERSION} for ${platform}...`);

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const follow = (url) => {
      https.get(url, (res) => {
        if (res.statusCode === 301 || res.statusCode === 302) {
          return follow(res.headers.location);
        }
        if (res.statusCode !== 200) {
          return reject(new Error(`Download failed: ${res.statusCode}`));
        }
        const file = fs.createWriteStream(dest);
        res.pipe(file);
        file.on('finish', () => {
          file.close();
          resolve();
        });
      }).on('error', reject);
    };
    follow(url);
  });
}

async function extract(archivePath, destDir) {
  if (archivePath.endsWith('.zip')) {
    // For Windows, use PowerShell
    execSync(`powershell -command "Expand-Archive -Path '${archivePath}' -DestinationPath '${destDir}' -Force"`, { stdio: 'inherit' });
  } else {
    // For tar.gz
    execSync(`tar -xzf "${archivePath}" -C "${destDir}"`, { stdio: 'inherit' });
  }
}

async function main() {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'cl-'));
  const archivePath = path.join(tempDir, archive);

  try {
    await download(url, archivePath);
    await extract(archivePath, binDir);

    // Make binary executable on Unix
    if (os.platform() !== 'win32') {
      fs.chmodSync(binPath, 0o755);
    }

    console.log(`cl v${VERSION} installed successfully!`);
  } catch (err) {
    console.error('Failed to install cl:', err.message);
    console.error('Please install manually from: https://github.com/groo-dev/cl-wranger/releases');
    process.exit(1);
  } finally {
    // Cleanup temp files
    fs.rmSync(tempDir, { recursive: true, force: true });
  }
}

main();

# cl - Cloudflare Wrangler Account Switcher

A CLI tool to easily switch between multiple Cloudflare/Wrangler accounts.

## Features

- **Save accounts** - Store multiple Wrangler authentication configs
- **Quick switch** - Switch between accounts with fuzzy matching
- **Interactive mode** - Select accounts from a list when you can't remember
- **Auto-save** - Automatically saves token updates before switching
- **Login flow** - Add new accounts directly from the switch menu

## Installation

### Homebrew (macOS/Linux)

```bash
brew install groo-dev/tap/cl
```

### npm

```bash
npm install -g @groo.dev/cl-wrangler
```

### pip

```bash
pip install cl-wrangler
```

### From Releases

Download the latest binary from [Releases](https://github.com/groo-dev/cl-wrangler/releases):

```bash
# macOS (Apple Silicon)
curl -L https://github.com/groo-dev/cl-wrangler/releases/latest/download/cl_darwin_arm64.tar.gz | tar xz
sudo mv cl /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/groo-dev/cl-wrangler/releases/latest/download/cl_darwin_amd64.tar.gz | tar xz
sudo mv cl /usr/local/bin/

# Linux (amd64)
curl -L https://github.com/groo-dev/cl-wrangler/releases/latest/download/cl_linux_amd64.tar.gz | tar xz
sudo mv cl /usr/local/bin/
```

### From Source

```bash
git clone https://github.com/groo-dev/cl-wrangler.git
cd cl-wrangler/cli
go build -o cl
sudo mv cl /usr/local/bin/
```

## Quick Start

```bash
# Save your current wrangler account
cl add

# Switch between accounts (interactive menu)
cl switch
```

## Usage

### `cl switch` - The main command

Run `cl switch` to open an interactive menu where you can:
- **Switch** to any saved account
- **Login** with a new account
- **Remove** an account you no longer need

```bash
cl switch
```

```
? Select account:
  > HamidRaza (hamid@example.com) [current]
    WorkAccount (work@company.com)
    ClientProject (client@example.com)
  ────────────────
    + Login with new account
    - Remove an account
```

You can also switch directly with fuzzy matching:

```bash
cl switch hamid    # Fuzzy match by name
cl switch work     # Partial match works too
```

### Other commands

| Command | Description |
|---------|-------------|
| `cl add` | Save current wrangler account |
| `cl list` | List all saved accounts |
| `cl current` | Show current account |
| `cl remove` | Remove an account (also available in `cl switch`) |
| `cl logout` | Logout and remove current account |
| `cl config` | View/edit configuration |
| `cl version` | Show version |

## Shell Completions

### Zsh

```bash
# Add to ~/.zshrc
eval "$(cl completion zsh)"

# Or generate to file
cl completion zsh > "${fpath[1]}/_cl"
```

### Bash

```bash
# Add to ~/.bashrc
eval "$(cl completion bash)"
```

### Fish

```bash
cl completion fish | source
```

## How it works

1. Wrangler stores authentication at `~/Library/Preferences/.wrangler/config/default.toml` (macOS)
2. `cl` saves copies of this file for each account in `~/Library/Application Support/cl-wrangler/`
3. When switching, `cl` copies the saved config back to Wrangler's location
4. Before switching, any token updates are saved automatically (detected via file hash)

## Configuration

Config is stored at:
- macOS: `~/Library/Application Support/cl-wrangler/accounts.json`
- Linux: `~/.config/cl-wrangler/accounts.json`

### Custom Wrangler command

If wrangler isn't in your PATH, `cl` will prompt you on first run. You can also set it manually:

```bash
cl config
# Then edit the wrangler command
```

Or set via environment variable:

```bash
export CL_WRANGLER_CMD="/path/to/wrangler"
```

## License

MIT License - see [LICENSE](LICENSE)

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

## Usage

### Save your current account

```bash
cl add
```

This saves your current Wrangler authentication and marks it as the active account.

### List saved accounts

```bash
cl list
```

Shows all saved accounts with the current one marked.

### Switch accounts

```bash
# Interactive mode - shows a list to choose from
cl switch

# Fuzzy match - switch by name or ID
cl switch hamid
cl switch work
```

### Add a new account

```bash
# From the switch menu, select "+ Login with new account"
cl switch

# Or use wrangler directly, then save
wrangler login
cl add
```

### Remove an account

```bash
cl remove        # Interactive selection
cl remove work   # Fuzzy match
```

### Logout

```bash
cl logout
```

Runs `wrangler logout` and removes the account from saved accounts.

### View/Edit config

```bash
cl config
```

### Check version

```bash
cl version
```

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

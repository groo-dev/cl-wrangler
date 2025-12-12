# Getting Started

`cl` is a CLI tool to easily switch between multiple Cloudflare/Wrangler accounts.

## Why?

If you work with multiple Cloudflare accounts (personal projects, client work, different organizations), you know the pain of logging in and out of Wrangler. `cl` solves this by:

- Saving your Wrangler authentication configs
- Letting you switch between them instantly
- Auto-saving token updates before switching

## Quick Start

```bash
# Install (choose one)
brew install groo-dev/tap/cl     # Homebrew
npm install -g @groo.dev/cl-wrangler  # npm
pip install cl-wrangler          # pip

# Save your current account
cl add

# Switch between accounts
cl switch
```

## How It Works

1. Wrangler stores authentication at `~/.wrangler/config/default.toml`
2. `cl` saves copies of this file for each account
3. When switching, `cl` copies the saved config back
4. Before switching, any token updates are saved automatically

That's it! No complex setup, no environment variables to manage.

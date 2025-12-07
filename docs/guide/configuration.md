# Configuration

## Config Location

Config is stored at:

| OS | Path |
|----|------|
| macOS | `~/Library/Application Support/cl-wrangler/accounts.json` |
| Linux | `~/.config/cl-wrangler/accounts.json` |

## Custom Wrangler Command

If `wrangler` isn't in your PATH, `cl` will prompt you on first run.

You can also set it manually:

```bash
cl config
# Then edit the wrangler command
```

Or via environment variable:

```bash
export CL_WRANGLER_CMD="/path/to/wrangler"
```

## Data Storage

`cl` stores:

1. **Account metadata** - Names, emails, account IDs in `accounts.json`
2. **Config files** - Copies of Wrangler's `default.toml` for each account

All data is stored locally. Nothing is sent to external servers.

## Wrangler Config Path

Wrangler stores authentication at:

| OS | Path |
|----|------|
| macOS | `~/Library/Preferences/.wrangler/config/default.toml` |
| Linux | `~/.config/.wrangler/config/default.toml` |

`cl` reads from and writes to this location when switching accounts.

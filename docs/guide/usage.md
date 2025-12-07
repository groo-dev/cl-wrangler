# Usage

## The Main Command: `cl switch`

`cl switch` is your one-stop command for managing accounts. Run it to:

- **Switch** to any saved account
- **Login** with a new account
- **Remove** an account you no longer need

```bash
cl switch
```

```
? Select account:
  > ✓ PersonalAccount (personal@example.com)
      WorkAccount (work@company.com)
      ClientProject (client@example.com)
    + Add new account
    × Remove account
```

### Fuzzy Matching

Skip the menu by typing part of the account name:

```bash
cl switch work     # Matches "WorkAccount"
cl switch personal # Matches "PersonalAccount"
cl switch client   # Matches "ClientProject"
```

## Other Commands

| Command | Description |
|---------|-------------|
| `cl add` | Save current wrangler account |
| `cl list` | List all saved accounts |
| `cl current` | Show current account |
| `cl remove` | Remove an account |
| `cl logout` | Logout and remove current account |
| `cl config` | View/edit configuration |
| `cl version` | Show version |

## Workflow Example

```bash
# Start with your personal account logged in
wrangler login
cl add  # Saves as "PersonalAccount"

# Add work account
wrangler login  # Login to work account
cl add  # Saves as "WorkAccount"

# Now switch between them
cl switch personal  # Switch to personal
cl switch work      # Switch to work

# Or use the interactive menu
cl switch
```

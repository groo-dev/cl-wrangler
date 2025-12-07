# cl

Go source code for the `cl` CLI tool.

See [cl-wrangler.groo.dev](https://cl-wrangler.groo.dev) for documentation.

## Development

```bash
# Build
go build -o cl

# Run
./cl version
./cl switch
```

## Project Structure

```
cli/
├── cmd/          # Cobra commands
├── internal/
│   ├── store/    # Account storage and config management
│   └── wrangler/ # Wrangler CLI integration
├── main.go       # Entry point
└── VERSION       # Current version
```

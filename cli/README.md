# cl

Go source code for the `cl` CLI tool.

See [docs.groo.dev/cl-wrangler](https://docs.groo.dev/cl-wrangler) for documentation.

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
│   ├── update/   # Version check
│   └── wrangler/ # Wrangler CLI integration
└── main.go       # Entry point
```

# Installation

Choose your preferred installation method:

## Homebrew (macOS/Linux)

```bash
brew install groo-dev/tap/cl
```

## npm

```bash
npm install -g @groo.dev/cl-wrangler
```

## pip

```bash
pip install cl-wrangler
```

## From Releases

Download the latest binary from [GitHub Releases](https://github.com/groo-dev/cl-wrangler/releases):

::: code-group

```bash [macOS (Apple Silicon)]
curl -L https://github.com/groo-dev/cl-wrangler/releases/latest/download/cl_darwin_arm64.tar.gz | tar xz
sudo mv cl /usr/local/bin/
```

```bash [macOS (Intel)]
curl -L https://github.com/groo-dev/cl-wrangler/releases/latest/download/cl_darwin_amd64.tar.gz | tar xz
sudo mv cl /usr/local/bin/
```

```bash [Linux (amd64)]
curl -L https://github.com/groo-dev/cl-wrangler/releases/latest/download/cl_linux_amd64.tar.gz | tar xz
sudo mv cl /usr/local/bin/
```

:::

## From Source

```bash
git clone https://github.com/groo-dev/cl-wrangler.git
cd cl-wrangler/cli
go build -o cl
sudo mv cl /usr/local/bin/
```

## Shell Completions

Enable tab completion for your shell:

::: code-group

```bash [Zsh]
# Add to ~/.zshrc
eval "$(cl completion zsh)"

# Or generate to file
cl completion zsh > "${fpath[1]}/_cl"
```

```bash [Bash]
# Add to ~/.bashrc
eval "$(cl completion bash)"
```

```bash [Fish]
cl completion fish | source
```

:::

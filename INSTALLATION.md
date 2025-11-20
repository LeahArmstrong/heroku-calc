# Installation Guide

## ✅ Successfully Published

The module is now available at:
- **GitHub**: https://github.com/leaharmstrong/heroku-calc
- **Module Path**: `github.com/leaharmstrong/heroku-calc`
- **Latest Version**: `v1.0.0`

## Installation Methods

### Method 1: Install Specific Version (Recommended)

```bash
go install github.com/leaharmstrong/heroku-calc@v1.0.0
```

This will:
- Download and compile the v1.0.0 release
- Install to `$(go env GOPATH)/bin/heroku-calc`
- Use the correct module path

### Method 2: Install Latest (After Proxy Refresh)

```bash
go install github.com/leaharmstrong/heroku-calc@latest
```

**Note**: If you get a module path conflict error, it means Go's module proxy still has cached data. Use Method 1 with the specific version tag, or wait a few minutes and try again.

### Method 3: Clone and Build

```bash
git clone https://github.com/leaharmstrong/heroku-calc.git
cd heroku-calc
go build -o heroku-calc
```

## Verifying Installation

Check the binary location:
```bash
which heroku-calc
# Should show: /Users/yourusername/go/bin/heroku-calc
```

Or find it manually:
```bash
ls -la $(go env GOPATH)/bin/heroku-calc
```

Test it works:
```bash
heroku-calc --help
```

## Adding to PATH

If `heroku-calc` command is not found, add Go's bin directory to your PATH:

### Bash/Zsh
Add to `~/.bashrc` or `~/.zshrc`:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then reload:
```bash
source ~/.bashrc  # or source ~/.zshrc
```

### Fish
```fish
set -Ua fish_user_paths (go env GOPATH)/bin
```

## Troubleshooting

### Error: "module declares its path as: github.com/egg/heroku-calc"

This happens if Go's module proxy cached an old version before we updated the module path.

**Solution 1**: Use the specific version tag
```bash
go install github.com/leaharmstrong/heroku-calc@v1.0.0
```

**Solution 2**: Clear Go's module cache
```bash
go clean -modcache
go install github.com/leaharmstrong/heroku-calc@latest
```

**Solution 3**: Wait for proxy refresh
The Go module proxy (proxy.golang.org) typically refreshes within a few minutes. Try again later.

### Error: "heroku-calc: command not found"

The binary is installed but not in your PATH.

**Solution**:
```bash
# Use full path
$(go env GOPATH)/bin/heroku-calc --help

# Or add to PATH (see "Adding to PATH" above)
```

### Error: "permission denied"

The binary doesn't have execute permissions.

**Solution**:
```bash
chmod +x $(go env GOPATH)/bin/heroku-calc
```

## Uninstalling

Remove the binary:
```bash
rm $(go env GOPATH)/bin/heroku-calc
```

Clean module cache:
```bash
go clean -modcache
```

## Upgrading

When a new version is released:

```bash
# Install specific version
go install github.com/leaharmstrong/heroku-calc@v1.1.0

# Or upgrade to latest
go install github.com/leaharmstrong/heroku-calc@latest
```

## Version Information

Check installed version:
```bash
heroku-calc --version  # (if version flag is added)
```

Or check the binary info:
```bash
ls -l $(go env GOPATH)/bin/heroku-calc
go version -m $(go env GOPATH)/bin/heroku-calc
```

## For Developers

Import in your Go projects:

```go
import (
    "github.com/leaharmstrong/heroku-calc/internal/config"
    "github.com/leaharmstrong/heroku-calc/internal/analysis"
)
```

Add to `go.mod`:
```go
require github.com/leaharmstrong/heroku-calc v1.0.0
```

## System Requirements

- **Go**: 1.23.2 or later
- **OS**: macOS, Linux, Windows
- **Heroku CLI**: Required for Heroku integration
- **Git**: Required for repository detection

## Post-Installation

After installing, you're ready to use the tool:

```bash
# From a Rails project directory
cd ~/my-rails-app
heroku-calc

# Or specify project path
heroku-calc --project ~/my-rails-app

# Or specify app name
heroku-calc --app my-heroku-app

# Get help
heroku-calc --help
```

See [README.md](README.md) for full usage documentation.

## Release Information

**v1.0.0** (2025-11-19)
- Initial release
- Complete BubbleTea TUI
- Database, Redis, and web tier analysis
- Performance recommendations
- Markdown report generation
- Multiple operation modes (dry-run, interactive, apply)

## Support

- **Issues**: https://github.com/leaharmstrong/heroku-calc/issues
- **Documentation**: See README.md, QUICKSTART.md
- **Source**: https://github.com/leaharmstrong/heroku-calc

# Development Guide

## Local Testing Without Releases

### Method 1: Build and Run Locally (Recommended)

Build the binary from your local source:

```bash
cd /Users/egg/Work/heroku-calc
go build -o heroku-calc

# Then run it from anywhere with --project flag
./heroku-calc --project ~/path/to/your/rails-app

# Or run it FROM the Rails project
cd ~/path/to/your/rails-app
/Users/egg/Work/heroku-calc/heroku-calc
```

### Method 2: Use `go run` (No Build Needed)

Run directly from source without building:

```bash
cd /Users/egg/Work/heroku-calc

# Run with project path
go run . --project ~/path/to/your/rails-app

# Or navigate to Rails project first
cd ~/path/to/your/rails-app
go run /Users/egg/Work/heroku-calc
```

### Method 3: Install from Local Source

Install to your GOPATH from local source:

```bash
cd /Users/egg/Work/heroku-calc
go install

# Now you can run from anywhere
cd ~/path/to/your/rails-app
heroku-calc
```

This installs the current local version to `$(go env GOPATH)/bin/heroku-calc`.

### Method 4: Create a Symlink (Convenient for Repeated Testing)

Create a symlink to your local binary:

```bash
cd /Users/egg/Work/heroku-calc
go build -o heroku-calc

# Create symlink in your PATH
ln -sf $(pwd)/heroku-calc /usr/local/bin/heroku-calc

# Or in GOPATH/bin
ln -sf $(pwd)/heroku-calc $(go env GOPATH)/bin/heroku-calc
```

Now `heroku-calc` will always use your local development version.

### Method 5: Test Different Rails Apps

Create a test script:

```bash
# Create test.sh in heroku-calc directory
cat > test.sh << 'EOF'
#!/bin/bash
set -e

# Build current version
go build -o heroku-calc

# Test with your Rails apps
echo "Testing with app 1..."
./heroku-calc --project ~/projects/rails-app-1 --dry-run

echo "Testing with app 2..."
./heroku-calc --project ~/projects/rails-app-2 --dry-run
EOF

chmod +x test.sh
./test.sh
```

## Quick Development Workflow

### Typical Development Cycle

```bash
# 1. Make code changes in your editor
vim internal/analysis/database.go

# 2. Build (fast, usually <5 seconds)
go build -o heroku-calc

# 3. Test immediately
./heroku-calc --project ~/my-rails-app

# 4. Iterate - repeat steps 1-3 as needed
```

### With Auto-Rebuild (Using `entr` or `watchexec`)

Install a file watcher (optional):

```bash
# Using homebrew on macOS
brew install entr

# Or watchexec
brew install watchexec
```

Auto-rebuild on file changes:

```bash
# Using entr
find . -name "*.go" | entr -r go build -o heroku-calc

# Using watchexec
watchexec -e go -r 'go build -o heroku-calc'
```

Then in another terminal, test your changes:

```bash
./heroku-calc --project ~/my-rails-app
```

## Testing Specific Features

### Test Without Full UI (Quick Validation)

Create test files to validate specific functionality:

```bash
# Create a test file
cat > test_analysis.go << 'EOF'
package main

import (
    "fmt"
    "github.com/leaharmstrong/heroku-calc/internal/heroku"
    "github.com/leaharmstrong/heroku-calc/internal/analysis"
    "github.com/leaharmstrong/heroku-calc/internal/pricing"
)

func main() {
    client, _ := heroku.NewClient("your-app-name")
    pricingData, _ := pricing.Get()

    analyzer := analysis.NewAnalyzer(client, pricingData)
    analyzer.LoadData()

    result, _ := analyzer.Analyze()

    fmt.Printf("Database Status: %s\n", result.DatabaseAnalysis.Status)
    fmt.Printf("Recommendations: %d\n", len(result.Recommendations))
}
EOF

# Run it
go run test_analysis.go
```

### Test with Different Flags

```bash
# Test read-only mode (default)
./heroku-calc --project ~/my-rails-app

# Test dry-run mode
./heroku-calc --project ~/my-rails-app --dry-run

# Test interactive mode
./heroku-calc --project ~/my-rails-app --interactive

# Test export
./heroku-calc --project ~/my-rails-app --export test-report.md
```

### Test with Different Apps

```bash
# Specify app name instead of using git detection
./heroku-calc --app my-staging-app
./heroku-calc --app my-production-app
```

## Debugging

### Enable Verbose Output

Add debug logging to your code:

```go
import "log"

// Add at the start of functions you want to debug
log.Printf("DEBUG: Loading data for app %s", appName)
log.Printf("DEBUG: Client initialized: %+v", client)
```

### Run with Go's Race Detector

Detect data races during testing:

```bash
go run -race . --project ~/my-rails-app
```

### Profile Performance

```bash
# CPU profile
go build -o heroku-calc
./heroku-calc --project ~/my-rails-app --cpuprofile=cpu.prof

# View profile
go tool pprof cpu.prof
```

## Testing Edge Cases

### Test with Missing Config

```bash
cd ~/my-rails-app
rm .heroku-calc.yml  # Remove config
/Users/egg/Work/heroku-calc/heroku-calc
```

### Test with Invalid Heroku App

```bash
./heroku-calc --app nonexistent-app-123456
```

### Test with No Git Remote

```bash
cd /tmp
mkdir test-app && cd test-app
/Users/egg/Work/heroku-calc/heroku-calc --app my-app
```

## Integration Testing

### Create a Test Rails App

```bash
# Create minimal test app structure
mkdir -p ~/test-rails-app/.git
cd ~/test-rails-app

# Add fake git remote
git init
git remote add heroku https://git.heroku.com/test-app.git

# Test
/Users/egg/Work/heroku-calc/heroku-calc
```

## Before Committing

### Pre-commit Checklist

```bash
# 1. Format code
go fmt ./...

# 2. Run tests (if you add them)
go test ./...

# 3. Check for issues
go vet ./...

# 4. Build successfully
go build -o heroku-calc

# 5. Test with real app
./heroku-calc --project ~/real-rails-app --dry-run

# 6. Verify help works
./heroku-calc --help
```

## Quick Reference

### Most Common Testing Commands

```bash
# Quick rebuild and test
go build -o heroku-calc && ./heroku-calc --project ~/my-app

# Test from Rails directory
cd ~/my-rails-app && /Users/egg/Work/heroku-calc/heroku-calc

# Test with go run (no build)
go run . --project ~/my-app

# Install local version globally
go install && heroku-calc --project ~/my-app
```

### Environment Variables for Testing

```bash
# Set custom Heroku CLI path
export HEROKU_CLI=/custom/path/to/heroku

# Enable debug mode (if implemented)
export DEBUG=1
```

## Tips

1. **Keep a Terminal Open**: Keep one terminal in the heroku-calc directory for quick rebuilds
2. **Use Short Aliases**: Create shell aliases for common commands:
   ```bash
   alias hc-build='cd /Users/egg/Work/heroku-calc && go build -o heroku-calc'
   alias hc-test='hc-build && ./heroku-calc --project ~/my-rails-app'
   ```
3. **Test Early, Test Often**: Build and test after each small change
4. **Use --dry-run**: Always test with `--dry-run` first to avoid accidental changes
5. **Check Multiple Apps**: Test with different Rails apps to catch edge cases

## Troubleshooting Local Development

### "Command Not Found" After `go install`

Make sure GOPATH/bin is in your PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Changes Not Reflected

Make sure you're running the newly built binary:
```bash
# Check which binary you're running
which heroku-calc

# Force rebuild
go clean && go build -o heroku-calc

# Or use full path
./heroku-calc --project ~/my-app
```

### Import Errors

If you see import errors after changes:
```bash
go mod tidy
go build -o heroku-calc
```

## Advanced: Test with Replace Directive

If you want to test changes in a separate Go project that imports heroku-calc:

```bash
# In your test project's go.mod
go mod edit -replace github.com/leaharmstrong/heroku-calc=/Users/egg/Work/heroku-calc
go mod tidy
```

This is useful if you're building tooling that uses heroku-calc as a library.

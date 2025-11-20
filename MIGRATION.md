# Module Path Migration

## Summary

The module path has been updated from `github.com/egg/heroku-calc` to `github.com/leaharmstrong/heroku-calc` to reflect the correct GitHub organization.

## Changes Made

### 1. Module Path (go.mod)
```diff
- module github.com/egg/heroku-calc
+ module github.com/leaharmstrong/heroku-calc
```

### 2. Import Statements (16 Go files)
All internal imports updated across:
- `main.go`
- `cmd/root.go`
- `internal/analysis/*.go` (5 files)
- `internal/heroku/*.go` (2 files)
- `internal/report/*.go` (1 file)
- `internal/ui/*.go` (7 files)

Example change:
```diff
- import "github.com/egg/heroku-calc/internal/config"
+ import "github.com/leaharmstrong/heroku-calc/internal/config"
```

### 3. Documentation (2 files)
- `README.md`: Updated repository references
- `GIT_SUMMARY.md`: Updated module path references

### 4. Git Remote
```bash
# Before
origin  git@github.com:LeahArmstrong/heroku-calc.git

# After
origin  https://github.com/leaharmstrong/heroku-calc.git
```

## Verification

✅ **Build Successful**
```bash
go build -o heroku-calc
# Builds without errors
```

✅ **Binary Works**
```bash
./heroku-calc --help
# Shows help correctly
```

✅ **No Old References**
```bash
grep -r "github.com/egg/heroku-calc" . --include="*.go" --include="*.md"
# No results found
```

## Git Commit

```
commit bfe0b1a
Author: Claude
Date:   2025-11-19

    refactor: update module path to leaharmstrong/heroku-calc

    - Change module path from github.com/egg/heroku-calc to github.com/leaharmstrong/heroku-calc
    - Update all import statements across 16 Go files
    - Update documentation references in README.md and GIT_SUMMARY.md
    - Verified build succeeds with new module path

    This prepares the repository for GitHub publication under leaharmstrong organization.
```

## Files Changed

| File | Changes |
|------|---------|
| `go.mod` | Module path |
| `main.go` | Import statement |
| `cmd/root.go` | Import statement |
| `internal/analysis/*.go` | Import statements (5 files) |
| `internal/heroku/*.go` | Import statements (2 files) |
| `internal/report/markdown.go` | Import statement |
| `internal/ui/*.go` | Import statements (7 files) |
| `README.md` | Repository references |
| `GIT_SUMMARY.md` | Module path references |

**Total**: 19 files, 33 insertions(+), 33 deletions(-)

## Next Steps

### 1. Push to GitHub
```bash
git push -u origin main
```

### 2. Create Release Tag
```bash
git tag -a v1.0.0 -m "Initial release: Heroku Config Analyzer"
git push origin v1.0.0
```

### 3. Update Go Module Cache (for users)
Users who want to install or update will use:
```bash
go install github.com/leaharmstrong/heroku-calc@latest
```

### 4. GitHub Repository Setup
On https://github.com/leaharmstrong/heroku-calc:

**Repository Settings**:
- Description: "BubbleTea TUI for analyzing and optimizing Heroku Rails application configurations"
- Website: (optional)
- Topics: `go`, `heroku`, `tui`, `bubbletea`, `rails`, `analyzer`, `cli`, `performance`

**Features to Enable**:
- ✅ Issues
- ✅ Discussions
- ✅ Wiki (optional)
- ✅ Projects (optional)

**Branch Protection** (recommended):
- Require pull request reviews
- Require status checks
- Require branches to be up to date

## Installation for Users

Once pushed, users can install with:

```bash
# Install latest version
go install github.com/leaharmstrong/heroku-calc@latest

# Or clone and build
git clone https://github.com/leaharmstrong/heroku-calc.git
cd heroku-calc
go build -o heroku-calc
```

## Import for Developers

Other Go projects can import:

```go
import (
    "github.com/leaharmstrong/heroku-calc/internal/config"
    "github.com/leaharmstrong/heroku-calc/internal/analysis"
)
```

## Migration Complete ✅

The repository is now properly configured for:
- GitHub hosting under leaharmstrong
- Go module discovery and installation
- Collaboration and contributions
- Public distribution

All references have been updated and verified working!

# Heroku Config Analyzer

A powerful BubbleTea TUI application for analyzing and optimizing Heroku Rails application configurations. Get performance recommendations, identify potential load balancing issues, and apply configuration changes directly from your terminal.

## Features

- **Automatic App Detection**: Detects Heroku app from git remotes
- **Comprehensive Analysis**:
  - Database connection pool sizing and capacity planning
  - Redis connection analysis for caching and Sidekiq
  - Web tier concurrency optimization (Puma workers and threads)
- **Performance Recommendations**: Get actionable suggestions with cost impact analysis
- **Interactive TUI**: Beautiful tab-based interface built with BubbleTea
- **Flexible Operation Modes**:
  - Read-only analysis
  - Dry-run (preview changes)
  - Interactive (apply changes with confirmation)
  - Batch apply mode
- **Markdown Reports**: Export detailed analysis reports
- **Up-to-date Pricing**: Uses current Heroku marketplace data with caching

## Installation

### Prerequisites

- Go 1.23 or later
- Heroku CLI installed and authenticated
- Git repository with Heroku remote

### Build from Source

```bash
git clone https://github.com/leaharmstrong/heroku-calc
cd heroku-calc
go build -o heroku-calc
```

### Install

```bash
go install github.com/leaharmstrong/heroku-calc@latest
```

## Usage

### Basic Usage

Run from within your Rails project directory:

```bash
heroku-calc
```

The tool will automatically detect your Heroku app from git remotes.

### Specify Project Path

```bash
heroku-calc --project /path/to/rails/app
```

### Specify App Name

```bash
heroku-calc --app my-heroku-app
```

### Operation Modes

**Read-Only Mode** (default):
```bash
heroku-calc
```

**Dry-Run Mode** (preview changes):
```bash
heroku-calc --dry-run
```

**Interactive Mode** (confirm each change):
```bash
heroku-calc --interactive
```

**Apply Mode** (auto-apply all recommended changes):
```bash
heroku-calc --apply  # Use with caution!
```

### Export Report

```bash
heroku-calc --export report.md
```

Or press `e` in the TUI to export.

## Configuration File

The tool creates a `.heroku-calc.yml` file in your project root to store safe environment variables and configuration:

```yaml
app_name: my-rails-app
git_remote: heroku
safe_env_vars:
  - DATABASE_URL
  - REDIS_URL
  - WEB_CONCURRENCY
  - RAILS_MAX_THREADS
  - SIDEKIQ_CONCURRENCY
excluded_env_vars:
  - SECRET_KEY_BASE
  - API_KEY
last_updated: 2025-11-19T10:00:00Z
```

## UI Navigation

### Keyboard Shortcuts

- `Tab` / `→`: Next tab
- `Shift+Tab` / `←`: Previous tab
- `↑` / `k`: Move cursor up
- `↓` / `j`: Move cursor down
- `Enter` / `Space`: Select/toggle item
- `a`: Apply selected actions (Actions tab only)
- `e`: Export markdown report
- `q` / `Ctrl+C`: Quit

### Tabs

1. **Overview**: Application summary and health status
2. **Env Vars**: Select which variables to track in `.heroku-calc.yml`
3. **Dynos**: View dyno formation and costs
4. **Addons**: List configured addons
5. **Analysis**: Detailed configuration analysis
6. **Actions**: Recommended changes with apply options

## Analysis Performed

### Database Connections

The tool analyzes your Postgres configuration:

- Calculates total connections required:
  - Web dynos: `quantity × WEB_CONCURRENCY × RAILS_MAX_THREADS`
  - Workers: `quantity × SIDEKIQ_CONCURRENCY`
- Compares against your plan's connection limit
- Recommends plan upgrades if buffer is <50%
- Identifies connection exhaustion risks

### Redis Configuration

Analyzes Redis/cache setup:

- Checks Redis connection limits vs. usage
- Validates `REDIS_POOL_SIZE` configuration
- Analyzes Sidekiq concurrency settings
- Recommends plan upgrades when utilization >80%

### Web Tier Optimization

Reviews Puma/web worker configuration:

- Calculates memory per thread
- Validates `WEB_CONCURRENCY` and `RAILS_MAX_THREADS`
- Identifies over-configuration (risk of R14 errors)
- Suggests optimal settings for dyno type

## Example Analysis Output

```
DATABASE CONNECTIONS - 🔴 Critical
  Plan: standard-0
  Max connections: 120
  Current usage: 4 dynos × 2 workers × 5 threads = 40 connections
  Sidekiq: 2 dynos × 10 threads = 20 connections
  Total required: 60 / 120 available (50.0% buffer)

  Issues:
  • Low buffer: 50.0% available (recommend 50%+ for bursts)

REDIS CONFIGURATION - 🟡 Warning
  Plan: premium-0
  Max connections: 40
  Missing: REDIS_POOL_SIZE not set

WEB TIER - 🟢 Optimal
  Dyno type: Performance-M (2560 MB RAM)
  WEB_CONCURRENCY: 2 workers
  RAILS_MAX_THREADS: 5 threads per worker
  Total threads: 10
  Memory per thread: 256 MB
```

## Recommendations

The tool provides actionable recommendations:

- **Critical**: Immediate action required (connection exhaustion risk)
- **High**: Important optimizations
- **Medium**: Configuration improvements
- **Low**: Optional enhancements

Each recommendation includes:
- Current vs. suggested configuration
- Cost impact (for plan upgrades)
- Whether it can be auto-applied
- Environment variable name (if applicable)

## Pricing Data

The tool uses a hybrid approach for pricing:

1. Attempts to load from cache (`~/.heroku-calc/cache/`)
2. Falls back to bundled pricing data
3. Caches fetched data for 24 hours

Bundled pricing includes:
- All Heroku dyno types
- Postgres plans (Mini through Premium)
- Redis plans (Mini through Premium)

## Safety Features

- **Read-only by default**: Won't modify anything without explicit flags
- **Dry-run mode**: Preview all changes before applying
- **Interactive mode**: Confirm each change individually
- **Auto-apply filtering**: Only applies changes marked as safe
- **Config file**: Prevents accidental exposure of secrets

## Limitations

- Requires Heroku CLI (API mode planned for future release)
- Currently supports Rails applications only
- Some recommendations require manual intervention (e.g., plan upgrades)
- Pricing data is current as of 2025-11-19

## Development

### Project Structure

```
heroku-calc/
├── cmd/                    # CLI commands
├── internal/
│   ├── analysis/           # Configuration analysis engine
│   ├── config/             # Config file management
│   ├── heroku/             # Heroku API/CLI client
│   ├── pricing/            # Pricing data management
│   ├── report/             # Markdown report generation
│   └── ui/                 # BubbleTea TUI
└── data/                   # Bundled pricing data
```

### Building

```bash
go build -o heroku-calc
```

### Testing

```bash
go test ./...
```

## Contributing

Contributions welcome! Please submit issues and pull requests on GitHub.

## License

MIT License - see LICENSE file for details

## Support

For issues and feature requests, please use the GitHub issue tracker.

## Credits

Built with:
- [BubbleTea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Cobra](https://github.com/spf13/cobra) - CLI framework

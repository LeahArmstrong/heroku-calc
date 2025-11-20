# Quick Start Guide

## Installation

```bash
cd /Users/egg/Work/heroku-calc
go build -o heroku-calc
```

## First Run

### Prerequisites

1. Ensure you have the Heroku CLI installed:
   ```bash
   heroku --version
   ```

2. Authenticate with Heroku:
   ```bash
   heroku login
   ```

3. Navigate to a Rails project with a Heroku git remote, or use the `--app` flag

### Run the Analyzer

From within your Rails project directory:

```bash
heroku-calc
```

Or from anywhere, specifying the project:

```bash
heroku-calc --project /path/to/your/rails/app
```

Or specify the Heroku app directly:

```bash
heroku-calc --app your-app-name
```

## Using the TUI

### Navigation

- **Tab** or **→**: Move to next tab
- **Shift+Tab** or **←**: Move to previous tab
- **↑/↓** or **k/j**: Navigate items
- **Enter** or **Space**: Select/toggle item
- **e**: Export markdown report
- **a**: Apply selected actions (on Actions tab, when not in read-only mode)
- **q** or **Ctrl+C**: Quit

### Tabs Overview

1. **Overview**: Quick summary of your app configuration and health status
2. **Env Vars**: Select which environment variables to track in `.heroku-calc.yml`
3. **Dynos**: View your dyno formation and estimated costs
4. **Addons**: List of configured addons
5. **Analysis**: Detailed analysis of database, Redis, and web tier
6. **Actions**: Recommended changes with the ability to apply them

## Operation Modes

### Read-Only (Default)

Just analyze, don't change anything:

```bash
heroku-calc
```

### Dry Run

See what would be changed without actually applying:

```bash
heroku-calc --dry-run
```

### Interactive

Confirm each change before applying:

```bash
heroku-calc --interactive
```

### Auto-Apply

Apply all auto-applicable recommendations:

```bash
heroku-calc --apply
```

**Warning**: Only use `--apply` if you understand the recommendations!

## Common Workflows

### First Time Setup

1. Run the analyzer:
   ```bash
   heroku-calc
   ```

2. On the **Env Vars** tab, select which variables you want to track
   - Use Space/Enter to toggle selection
   - The tool automatically saves to `.heroku-calc.yml`

3. Navigate to the **Analysis** tab to see detailed findings

4. Check the **Actions** tab for recommendations

5. Export a report:
   - Press **e** in the TUI, or
   - Run with `--export` flag

### Regular Health Checks

```bash
# Quick check
heroku-calc

# Navigate to Analysis tab to see current status
# Green (🟢) = Optimal
# Yellow (🟡) = Warning
# Red (🔴) = Critical
```

### Applying Configuration Changes

1. Run in interactive mode:
   ```bash
   heroku-calc --interactive
   ```

2. Navigate to the **Actions** tab

3. Review recommendations (use ↑/↓ to see details)

4. Toggle which actions to apply (Space/Enter)

5. Press **a** to apply selected actions

6. Confirm in the prompts

### Exporting Reports

Export to default filename:

```bash
heroku-calc --export
```

Or specify filename:

```bash
heroku-calc --export analysis-2025-11-19.md
```

Or press **e** while in the TUI.

## Understanding the Analysis

### Database Connections

The tool calculates:
- **Web connections**: `web_dynos × WEB_CONCURRENCY × RAILS_MAX_THREADS`
- **Worker connections**: `worker_dynos × SIDEKIQ_CONCURRENCY`
- **Total vs. Max**: Compares your needs against your plan's limits
- **Buffer %**: Remaining capacity for bursts

**Critical**: <20% buffer (upgrade plan immediately)
**Warning**: <50% buffer (consider upgrading)
**Optimal**: ≥50% buffer

### Redis Configuration

Analyzes:
- Connection capacity vs. usage
- `REDIS_POOL_SIZE` configuration
- Sidekiq concurrency settings

### Web Tier

Evaluates:
- Memory per thread (recommended: 80+ MB)
- Thread count appropriateness for dyno size
- `WEB_CONCURRENCY` and `RAILS_MAX_THREADS` settings

## Configuration File

The tool creates `.heroku-calc.yml` in your project:

```yaml
app_name: my-rails-app
git_remote: heroku
safe_env_vars:
  - DATABASE_URL
  - REDIS_URL
  - WEB_CONCURRENCY
  - RAILS_MAX_THREADS
excluded_env_vars:
  - SECRET_KEY_BASE
last_updated: 2025-11-19T20:00:00Z
```

**Tip**: Add `.heroku-calc.yml` to your git repository so your team can use the same configuration.

## Troubleshooting

### "Failed to detect Heroku app"

- Ensure you're in a git repository
- Check you have a Heroku remote: `git remote -v`
- Or use `--app` flag: `heroku-calc --app your-app-name`

### "Failed to connect to Heroku"

- Ensure Heroku CLI is installed: `heroku --version`
- Login to Heroku: `heroku login`
- Verify access to the app: `heroku apps:info -a your-app-name`

### "No analysis available"

- The tool may still be loading data
- Wait for the spinner to complete
- Check for error messages in the status bar

### Build Errors

If you need to rebuild:

```bash
go mod tidy
go build -o heroku-calc
```

## Tips & Best Practices

1. **Run regularly**: Check your config after deploying major changes
2. **Review before applying**: Always review recommendations in dry-run mode first
3. **Start conservative**: Begin with read-only mode, move to interactive/apply once comfortable
4. **Track safe vars only**: Don't include secrets in `.heroku-calc.yml`
5. **Export reports**: Keep reports for capacity planning discussions
6. **Monitor trends**: Compare reports over time to see how your app evolves

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Explore all tabs in the TUI to familiarize yourself with the interface
- Run in different modes to understand the workflow
- Export a report and share with your team

## Getting Help

If you encounter issues:
1. Check the error message in the status bar (bottom of screen)
2. Verify your Heroku CLI is working: `heroku apps:info`
3. Review the README.md for detailed documentation
4. File an issue on GitHub with details about the error

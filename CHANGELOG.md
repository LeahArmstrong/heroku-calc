# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2025-11-20

### Fixed
- **Critical**: Fixed nil pointer dereference when loading Heroku data
- Heroku client now properly passed from data loading phase to analysis phase
- Resolved panic that occurred during environment variable loading

### Technical Details
The `loadData` function was creating a Heroku client but not including it in the returned message. When the analysis phase tried to use the client, it encountered a nil pointer, causing the application to crash. This has been fixed by:
- Adding `client` field to `loadedDataMsg` struct
- Including the client in the returned message
- Properly setting `m.herokuClient` when processing the loaded data message

## [1.0.0] - 2025-11-19 [DEPRECATED]

### Added
- Initial release of Heroku Config Analyzer
- BubbleTea TUI with tab-based dashboard (Overview, Env Vars, Dynos, Addons, Analysis, Actions)
- Comprehensive analysis engine for:
  - Database connection pool sizing and capacity planning
  - Redis connection analysis for caching and Sidekiq
  - Web tier concurrency optimization (Puma workers and threads)
- Performance recommendations with severity levels (critical, high, medium, low)
- Flexible operation modes: read-only, dry-run, interactive, apply
- Markdown report generation with detailed tables
- Heroku CLI integration with fallback support
- Git remote auto-detection
- YAML-based configuration management (`.heroku-calc.yml`)
- Up-to-date Heroku marketplace pricing data with caching
- Cost impact analysis for plan upgrades
- Auto-apply capabilities for safe configuration changes
- Comprehensive documentation (README, QUICKSTART, examples)

### Known Issues
- Nil pointer dereference when loading Heroku data (fixed in v1.0.1)

### Deprecated
- v1.0.0 is deprecated due to critical bug. Use v1.0.1 or later.

## Installation

### Latest Stable
```bash
go install github.com/leaharmstrong/heroku-calc@latest
```

### Specific Version
```bash
go install github.com/leaharmstrong/heroku-calc@v1.0.1
```

## Version Support

| Version | Status | Support | Notes |
|---------|--------|---------|-------|
| 1.0.1 | ✅ Stable | Full support | Current stable release |
| 1.0.0 | ⚠️ Deprecated | None | Critical bug - do not use |

## Upgrade Guide

### From v1.0.0 to v1.0.1

Simply reinstall with the new version:

```bash
go install github.com/leaharmstrong/heroku-calc@v1.0.1
```

No configuration changes required. The fix is transparent to users.

## Future Roadmap

Planned features for future releases:

### v1.1.0 (Planned)
- Direct Heroku API implementation (no CLI dependency)
- Support for Node.js applications
- Historical trend analysis
- Custom threshold configuration

### v1.2.0 (Planned)
- GitHub Actions integration
- Slack/email notifications for critical issues
- Multi-app comparison
- Cost optimization reports

### v2.0.0 (Future)
- Support for AWS, GCP, Azure
- Plugin system for custom analyzers
- Web UI option
- Real-time monitoring

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:
- Reporting bugs
- Suggesting features
- Submitting pull requests
- Development setup

## Links

- **Repository**: https://github.com/leaharmstrong/heroku-calc
- **Issues**: https://github.com/leaharmstrong/heroku-calc/issues
- **Releases**: https://github.com/leaharmstrong/heroku-calc/releases
- **Documentation**: See README.md and QUICKSTART.md

[1.0.1]: https://github.com/leaharmstrong/heroku-calc/releases/tag/v1.0.1
[1.0.0]: https://github.com/leaharmstrong/heroku-calc/releases/tag/v1.0.0

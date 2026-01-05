# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-01-05

### Added
- Initial release with core functionality
- UK tax and national insurance calculations via listentotaxman.com API
- Support for multiple regions (England/UK, Scotland, Wales, Northern Ireland)
- Student loan repayment calculations (Plan 1, 2, 4, Postgraduate, Scottish)
- Marriage allowance support with partner income
- Blind person's allowance
- National Insurance exemption option
- Pension contribution calculations (percentage or fixed amount)
- Period conversion (yearly, monthly, weekly, daily, hourly)
- Compare mode for side-by-side scenario comparison (2-4 scenarios)
- JSON output support for scripting
- Verbose mode with detailed tax breakdowns
- Configuration file support at `~/.config/listentotaxman/config.yaml`
- Shell completions (bash, zsh, fish, powershell)
- Homebrew tap support (`brew install mheap/tap/listentotaxman`)
- Docker images on Docker Hub (`mheap/listentotaxman`)
- Pre-built binaries for Linux, macOS, and Windows (amd64, arm64)

### Changed
- Renamed `--grosswage` flag to `--income` (breaking change from conceptual v0.0.x)
- Updated date logic to use April 5th cutoff (was April 3rd in initial development)

[Unreleased]: https://github.com/mheap/listentotaxman-cli/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/mheap/listentotaxman-cli/releases/tag/v0.1.0

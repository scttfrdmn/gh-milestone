# gh-milestone

A GitHub CLI extension for comprehensive milestone management.

## Why This Extension?

The GitHub CLI doesn't provide native milestone commands, requiring users to use raw REST API calls for all milestone operations. This extension provides intuitive, full-featured milestone management.

## Installation

```bash
gh extension install scttfrdmn/gh-milestone
```

## Commands

### Create Milestone

```bash
gh milestone create "v1.0.0" \
  --title "Version 1.0.0 Release" \
  --description "First stable release" \
  --due-date "2025-12-31"
```

**Flags:**
- `--title` / `-t` (required): Milestone title
- `--description` / `-d`: Description
- `--due-date`: Due date (YYYY-MM-DD format)
- `--state`: `open` (default) or `closed`
- `--repo` / `-R`: Target repository (default: current repo)

### List Milestones

```bash
gh milestone list
gh milestone list --state closed
gh milestone list --state all
```

**Flags:**
- `--state` / `-s`: Filter by state (`open`, `closed`, `all`) - default: `open`
- `--sort`: Sort by (`due-on`, `completeness`, `title`) - default: `due-on`
- `--json`: Output JSON
- `--repo` / `-R`: Target repository

### View Milestone

```bash
gh milestone view 1
gh milestone view "v1.0.0"
```

View detailed information about a milestone including progress and issue list.

### Edit Milestone

```bash
gh milestone edit 1 --title "v1.0.0 - Updated Title"
gh milestone edit 1 --due-date "2026-01-15"
gh milestone edit 1 --state closed
```

**Flags:**
- `--title` / `-t`: Update title
- `--description` / `-d`: Update description
- `--due-date`: Update due date
- `--state`: Update state (`open` or `closed`)

### Close/Reopen Milestone

```bash
gh milestone close 1
gh milestone reopen 1
```

Convenient shortcuts for changing milestone state.

### Delete Milestone

```bash
gh milestone delete 1
gh milestone delete 1 --yes  # Skip confirmation
```

## Development

### Prerequisites

- Go 1.21 or later
- GitHub CLI (`gh`) installed

### Local Installation

```bash
git clone https://github.com/scttfrdmn/gh-milestone.git
cd gh-milestone
go build
gh extension install .
```

### Testing

```bash
# Test the extension locally
gh milestone list

# Run Go tests
go test ./...
```

### Building

```bash
# Build for current platform
go build -o gh-milestone

# Cross-compile for release (done automatically by GitHub Actions)
GOOS=linux GOARCH=amd64 go build -o gh-milestone-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o gh-milestone-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o gh-milestone-darwin-arm64
GOOS=windows GOARCH=amd64 go build -o gh-milestone-windows-amd64.exe
```

## Architecture

```
gh-milestone/
├── main.go              # Entry point and CLI setup
├── cmd/                 # Command implementations
│   ├── create.go
│   ├── list.go
│   ├── view.go
│   ├── edit.go
│   └── delete.go
├── pkg/
│   ├── api/            # GitHub API client wrapper
│   └── format/         # Output formatting (tables, JSON)
└── .github/
    └── workflows/
        └── release.yml  # Automated cross-platform builds
```

## Release Process

1. Update version in code
2. Commit changes: `git commit -am "Release v1.0.0"`
3. Create tag: `git tag v1.0.0`
4. Push tag: `git push --tags`
5. GitHub Actions automatically builds cross-platform binaries and creates release

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

## License

MIT License - see LICENSE file for details

## Related

- [gh-label-sync](https://github.com/scttfrdmn/gh-label-sync) - Companion extension for label management
- [GitHub CLI](https://cli.github.com/)
- [go-gh](https://github.com/cli/go-gh) - Go library for building GitHub CLI extensions

## Background

This extension was created because GitHub CLI maintainers [decided not to include](https://github.com/cli/cli/issues/1200) native milestone commands in the core CLI, preferring users create extensions for this functionality.

The design is based on comprehensive [real-world experience](../research/cargoship-setup-feedback.md) setting up GitHub projects.

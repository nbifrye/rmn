# rmn

[![CI](https://github.com/nbifrye/rmn/actions/workflows/ci.yml/badge.svg)](https://github.com/nbifrye/rmn/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/nbifrye/rmn)](https://github.com/nbifrye/rmn/blob/main/go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/nbifrye/rmn/blob/main/LICENSE)
[![Docs](https://img.shields.io/badge/docs-GitHub%20Pages-blue)](https://nbifrye.github.io/rmn/)

**[日本語版 README はこちら](README.ja.md)**

**rmn** is an unofficial command-line client for [Redmine](https://www.redmine.org/) written in Go. It provides a fast, intuitive interface for managing Redmine issues, projects, users, versions, time entries, memberships, wiki pages, and more, directly from your terminal. Inspired by [GitLab CLI (glab)](https://gitlab.com/gitlab-org/cli), rmn brings familiar command patterns to the Redmine ecosystem.

rmn also includes a built-in [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server, enabling AI agents such as Claude Code to interact with your Redmine instance through natural language. Whether you prefer the keyboard or an AI assistant, rmn gives you full control over Redmine issue management via the Redmine REST API.

> **Note:** This project is not affiliated with or endorsed by the Redmine project. It is an independent, community-driven tool.

> **Warning:** This project is experimental. The API and CLI interface are not yet stable and may introduce breaking changes without notice.

## Features

- **Broad Redmine API coverage** -- issues, projects, users, versions, time entries, memberships, wiki pages, trackers and issue statuses
- **Full issue lifecycle management** -- list, view, create, update, close, and delete Redmine issues from the command line
- **MCP server for AI agents** -- expose Redmine operations to AI assistants like Claude Code via the Model Context Protocol
- **Multiple output formats** -- human-readable table (default) and machine-readable JSON for scripting and automation
- **Flexible issue filtering** -- filter by project, status, assignee, and tracker with sorting and pagination support
- **GitLab CLI-inspired aliases** -- familiar shorthand commands (`ls`, `show`, `get`, `new`, `rm`) for faster workflows
- **Six installation methods** -- Homebrew, mise, Nix, Go install, pre-built binaries, and build from source
- **Shell completion** -- auto-complete for Bash, Zsh, Fish, and PowerShell
- **XDG-compliant configuration** -- respects `$XDG_CONFIG_HOME` for config file placement
- **Security hardened** -- TLS 1.2+ enforcement, secure file permissions (0600), and API key protection on redirects
- **100% test coverage** -- enforced in CI with every pull request
- **Cross-platform** -- pre-built binaries for Linux, macOS, and Windows on both amd64 and arm64

## Quick Start

```bash
# 1. Install
brew tap nbifrye/rmn https://github.com/nbifrye/rmn.git
brew install nbifrye/rmn/rmn

# 2. Authenticate
rmn auth login --url https://your-redmine.example.com --api-key YOUR_API_KEY

# 3. List your issues
rmn issue list -a me

# 4. View issue details
rmn issue view 42
```

## Installation

Pre-built binaries are available for Linux, macOS, and Windows on both amd64 and arm64 architectures.

### Homebrew (macOS/Linux)

```bash
brew tap nbifrye/rmn https://github.com/nbifrye/rmn.git
brew install nbifrye/rmn/rmn
```

### mise

```bash
mise use -g ubi:nbifrye/rmn
```

### Nix

```bash
nix profile install github:nbifrye/rmn
```

### Go

```bash
go install github.com/nbifrye/rmn/cmd/rmn@latest
```

### Binary download

Download pre-built binaries from [GitHub Releases](https://github.com/nbifrye/rmn/releases).

### Build from source

Requires Go 1.24 or later.

```bash
make build
```

## Configuration

### Authentication

```bash
# Set up with flags
rmn auth login --url https://your-redmine.example.com --api-key YOUR_API_KEY

# Or interactively
rmn auth login

# Verify your configuration
rmn auth status
```

### Config file

Configuration is stored in `~/.config/rmn/config.json` (or `$XDG_CONFIG_HOME/rmn/config.json` if set):

```json
{
  "redmine_url": "https://your-redmine.example.com",
  "api_key": "your-api-key-here"
}
```

The config file is created with `0600` permissions. rmn refuses to read config files with insecure permissions.

## Documentation

For detailed documentation, visit the [rmn documentation site](https://nbifrye.github.io/rmn/).

| Topic | Description |
|-------|-------------|
| [Usage Guide](https://nbifrye.github.io/rmn/guide/usage) | Commands for issues, projects, users, versions, time entries, memberships, wiki pages, and more |
| [MCP Server](https://nbifrye.github.io/rmn/mcp-server) | Built-in Model Context Protocol server for AI agent integration (35 tools) |
| [Configuration](https://nbifrye.github.io/rmn/guide/configuration) | Full configuration reference including per-command overrides |
| [Shell Completion](https://nbifrye.github.io/rmn/reference/shell-completion) | Auto-complete setup for Bash, Zsh, Fish, and PowerShell |
| [Security](https://nbifrye.github.io/rmn/reference/security) | TLS enforcement, file permissions, API key protection |
| [Architecture](https://nbifrye.github.io/rmn/reference/architecture) | Project structure and key dependencies |

## Development

```bash
make build    # Build binary
make test     # Run all tests
make vet      # Static analysis
make lint     # Run linters (requires golangci-lint)
make cover    # Coverage report (enforces 100% coverage)
make install  # Install to $GOPATH/bin
make clean    # Remove build artifacts
```

All pull requests must pass CI, which enforces:
- `go vet ./...` with zero warnings
- 100% test coverage

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. **Update documentation** if your change affects user-facing behavior (see `CLAUDE.md` for the code → docs mapping)
5. Ensure tests pass: `make test && make vet`
6. Commit and push
7. Open a pull request

Please follow the existing code style: one file per subcommand, table-driven tests with `httptest.NewServer`, and `fmt.Errorf("context: %w", err)` for error wrapping.

## License

This project is licensed under the [MIT License](LICENSE).
